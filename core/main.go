package main

import (
	"path/filepath"
	"strconv"
	"sync"

	"github.com/mesg-foundation/engine/config"
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/database"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/logger"
	"github.com/mesg-foundation/engine/orchestrator"
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

func initDatabases(cfg *config.Config) (*database.LevelDBInstanceDB, *database.LevelDBExecutionDB, *database.LevelDBProcessDB, error) {
	// init instance db.
	instanceDB, err := database.NewInstanceDB(filepath.Join(cfg.Path, cfg.Database.InstanceRelativePath))
	if err != nil {
		return nil, nil, nil, err
	}

	// init execution db.
	executionDB, err := database.NewExecutionDB(filepath.Join(cfg.Path, cfg.Database.ExecutionRelativePath))
	if err != nil {
		return nil, nil, nil, err
	}

	// init process db.
	processDB, err := database.NewProcessDB(filepath.Join(cfg.Path, cfg.Database.ProcessRelativePath))
	if err != nil {
		return nil, nil, nil, err
	}

	return instanceDB, executionDB, processDB, nil
}

func deployCoreServices(cfg *config.Config, sdk *enginesdk.SDK) error {
	for _, serviceConfig := range cfg.SystemServices {
		logrus.WithField("module", "main").Infof("Deploying service %q", serviceConfig.Definition.Sid)
		srv, err := sdk.Service.Create(serviceConfig.Definition, "", "")
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
		logrus.WithField("module", "main").Infof("Service %q deployed with hash %q", srv.Sid, srv.Hash.String())
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
	cfg, err := config.Global()
	if err != nil {
		logrus.Fatalln(err)
	}

	// init logger.
	logger.Init(cfg.Log.Format, cfg.Log.Level, cfg.Log.ForceColors)

	// init databases
	instanceDB, executionDB, processDB, err := initDatabases(cfg)
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	// init container.
	c, err := container.New(cfg.Name)
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	_, port, _ := xnet.SplitHostPort(cfg.Server.Address)

	// init app factory
	db, err := db.NewGoLevelDB("app", cfg.Cosmos.Path)
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}
	appFactory := cosmos.NewAppFactory(logger.TendermintLogger(), db)

	// register the backend modules to the app factory.
	enginesdk.NewBackend(appFactory, c)

	// init cosmos app
	app, err := cosmos.NewApp(appFactory)
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	// init key manager
	kb, err := cosmos.NewKeybase(cfg.Cosmos.Path)
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	// create cosmos node
	node, err := cosmos.NewNode(app, kb, cfg.Tendermint.Config, &cfg.Cosmos)
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	// create cosmos client
	client := cosmos.NewClient(node, app.Cdc(), kb, cfg.Cosmos.ChainID)

	// init sdk
	sdk := enginesdk.New(client, app.Cdc(), kb, c, instanceDB, executionDB, processDB, cfg.Name, strconv.Itoa(port))

	// start tendermint node
	logrus.WithField("module", "main").WithField("seeds", cfg.Tendermint.P2P.Seeds).Info("starting tendermint node")
	if err := node.Start(); err != nil {
		logrus.WithField("module", "main").Fatalln(err)
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

	logrus.WithField("module", "main").Info("starting process engine")
	s := orchestrator.New(sdk.Event, sdk.Execution, sdk.Process)
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
