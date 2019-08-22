package main

import (
	"flag"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/mesg-foundation/engine/config"
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/database"
	"github.com/mesg-foundation/engine/database/store"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/logger"
	"github.com/mesg-foundation/engine/scheduler"
	enginesdk "github.com/mesg-foundation/engine/sdk"
	instancesdk "github.com/mesg-foundation/engine/sdk/instance"
	servicesdk "github.com/mesg-foundation/engine/sdk/service"
	"github.com/mesg-foundation/engine/server/grpc"
	"github.com/mesg-foundation/engine/version"
	"github.com/mesg-foundation/engine/x/xerrors"
	"github.com/mesg-foundation/engine/x/xnet"
	"github.com/mesg-foundation/engine/x/xos"
	"github.com/mesg-foundation/engine/x/xsignal"
	"github.com/sirupsen/logrus"
	db "github.com/tendermint/tm-db"
)

var network = flag.Bool("experimental-network", false, "start the engine with the network")

func initDatabases(cfg *config.Config) (*database.ServiceDB, *database.LevelDBInstanceDB, *database.LevelDBExecutionDB, *database.LevelDBWorkflowDB, error) {
	// init services db.
	serviceStore, err := store.NewLevelDBStore(filepath.Join(cfg.Path, cfg.Database.ServiceRelativePath))
	if err != nil {
		return nil, nil, nil, nil, err
	}
	serviceDB := database.NewServiceDB(serviceStore)

	// init instance db.
	instanceDB, err := database.NewInstanceDB(filepath.Join(cfg.Path, cfg.Database.InstanceRelativePath))
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// init execution db.
	executionDB, err := database.NewExecutionDB(filepath.Join(cfg.Path, cfg.Database.ExecutionRelativePath))
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// init workflow db.
	workflowDB, err := database.NewWorkflowDB(filepath.Join(cfg.Path, cfg.Database.WorkflowRelativePath))
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return serviceDB, instanceDB, executionDB, workflowDB, nil
}

func deployCoreServices(cfg *config.Config, sdk *enginesdk.SDK) error {
	for _, serviceConfig := range cfg.SystemServices {
		logrus.WithField("module", "main").Infof("Deploying service %q", serviceConfig.Definition.Sid)
		srv, err := sdk.Service.Create(serviceConfig.Definition)
		if err != nil {
			existsError, ok := err.(*servicesdk.AlreadyExistsError)
			if ok {
				srv, err = sdk.Service.Get(existsError.Hash)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
		logrus.WithField("module", "main").Infof("Service %q deployed with hash %q", srv.Sid, srv.Hash)
		instance, err := sdk.Instance.Create(srv.Hash, xos.EnvMapToSlice(serviceConfig.Env))
		if err != nil {
			existsError, ok := err.(*instancesdk.AlreadyExistsError)
			if ok {
				instance, err = sdk.Instance.Get(existsError.Hash)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
		serviceConfig.Instance = instance
		logrus.WithField("module", "main").Infof("Instance started with hash %q", instance.Hash)
	}
	return nil
}

func stopRunningServices(sdk *enginesdk.SDK) error {
	instances, err := sdk.Instance.List(&instancesdk.Filter{})
	if err != nil {
		return err
	}
	var (
		instancesLen = len(instances)
		errC         = make(chan error, instancesLen)
		wg           sync.WaitGroup
	)
	wg.Add(instancesLen)
	for _, instance := range instances {
		go func(hash hash.Hash) {
			defer wg.Done()
			err := sdk.Instance.Delete(hash, false)
			if err != nil {
				errC <- err
			}
		}(instance.Hash)
	}
	wg.Wait()
	close(errC)
	var errs xerrors.Errors
	for err := range errC {
		errs = append(errs, err)
	}
	return errs.ErrorOrNil()
}

func main() {
	flag.Parse()
	cfg, err := config.Global()
	if err != nil {
		logrus.Fatalln(err)
	}

	// init logger.
	logger.Init(cfg.Log.Format, cfg.Log.Level, cfg.Log.ForceColors)

	// init databases
	serviceDB, instanceDB, executionDB, workflowDB, err := initDatabases(cfg)
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	// init container.
	c, err := container.New(cfg.Name)
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	_, port, _ := xnet.SplitHostPort(cfg.Server.Address)

	var sdk *enginesdk.SDK
	if *network {
		// init cosmos app
		db, err := db.NewGoLevelDB("app", cfg.Cosmos.Path)
		if err != nil {
			logrus.WithField("module", "main").Fatalln(err)
		}
		app := cosmos.NewApp(logger.TendermintLogger(), db)

		// init sdk.
		sdk, err = enginesdk.New(app, c, serviceDB, instanceDB, executionDB, workflowDB, cfg.Name, strconv.Itoa(port))
		if err != nil {
			logrus.WithField("module", "main").Fatalln(err)
		}

		// create tendermint node
		node, err := cosmos.NewNode(app, cfg.Tendermint.Config, &cfg.Cosmos)
		if err != nil {
			logrus.WithField("module", "main").Fatalln(err)
		}

		// start tendermint node
		logrus.WithField("module", "main").WithField("seeds", cfg.Tendermint.P2P.Seeds).Info("starting tendermint node")
		if err := node.Start(); err != nil {
			logrus.WithField("module", "main").Fatalln(err)
		}
	} else {
		sdk = enginesdk.NewDeprecated(c, serviceDB, instanceDB, executionDB, workflowDB, cfg.Name, strconv.Itoa(port))
	}

	// init system services.
	if err := deployCoreServices(cfg, sdk); err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	// init gRPC server.
	server := grpc.New(sdk)

	logrus.WithField("module", "main").Infof("starting MESG Engine version %s", version.Version)

	go func() {
		if err := server.Serve(cfg.Server.Address); err != nil {
			logrus.WithField("module", "main").Fatalln(err)
		}
	}()

	logrus.WithField("module", "main").Info("starting workflow engine")
	s := scheduler.New(sdk.Event, sdk.Execution, sdk.Workflow)
	go func() {
		if err := s.Start(); err != nil {
			logrus.WithField("module", "main").Fatalln(err)
		}
	}()
	go func() {
		for err := range s.ErrC {
			logrus.WithField("module", "main").Warn(err)
		}
	}()

	<-xsignal.WaitForInterrupt()
	if err := stopRunningServices(sdk); err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}
	server.Close()
}
