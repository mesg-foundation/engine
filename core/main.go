package main

import (
	"flag"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/mesg-foundation/engine/config"
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/database"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/logger"
	"github.com/mesg-foundation/engine/scheduler"
	"github.com/mesg-foundation/engine/sdk"
	instancesdk "github.com/mesg-foundation/engine/sdk/instance"
	servicesdk "github.com/mesg-foundation/engine/sdk/service"
	"github.com/mesg-foundation/engine/server/grpc"
	"github.com/mesg-foundation/engine/tendermint"
	"github.com/mesg-foundation/engine/version"
	"github.com/mesg-foundation/engine/x/xerrors"
	"github.com/mesg-foundation/engine/x/xnet"
	"github.com/mesg-foundation/engine/x/xos"
	"github.com/mesg-foundation/engine/x/xsignal"
	"github.com/sirupsen/logrus"
)

var network = flag.Bool("experimental-network", false, "start the engine with the network")

type dependencies struct {
	cfg         *config.Config
	serviceDB   database.ServiceDB
	executionDB database.ExecutionDB
	workflowDB  database.WorkflowDB
	container   container.Container
	sdk         *sdk.SDK
}

func initDependencies(cfg *config.Config) (*dependencies, error) {
	// init services db.
	serviceDB, err := database.NewServiceDB(filepath.Join(cfg.Path, cfg.Database.ServiceRelativePath))
	if err != nil {
		return nil, err
	}

	// init instance db.
	instanceDB, err := database.NewInstanceDB(filepath.Join(cfg.Path, cfg.Database.InstanceRelativePath))
	if err != nil {
		return nil, err
	}

	// init execution db.
	executionDB, err := database.NewExecutionDB(filepath.Join(cfg.Path, cfg.Database.ExecutionRelativePath))
	if err != nil {
		return nil, err
	}

	// init workflow db.
	workflowDB, err := database.NewWorkflowDB(filepath.Join(cfg.Path, cfg.Database.WorkflowRelativePath))
	if err != nil {
		return nil, err
	}

	// init container.
	c, err := container.New(cfg.Name)
	if err != nil {
		return nil, err
	}

	_, port, _ := xnet.SplitHostPort(cfg.Server.Address)

	// init sdk.
	sdk := sdk.New(c, serviceDB, instanceDB, executionDB, workflowDB, cfg.Name, strconv.Itoa(port))

	return &dependencies{
		cfg:         cfg,
		container:   c,
		serviceDB:   serviceDB,
		executionDB: executionDB,
		workflowDB:  workflowDB,
		sdk:         sdk,
	}, nil
}

func deployCoreServices(cfg *config.Config, sdk *sdk.SDK) error {
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

func stopRunningServices(sdk *sdk.SDK) error {
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

	dep, err := initDependencies(cfg)
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	// init logger.
	logger.Init(cfg.Log.Format, cfg.Log.Level, cfg.Log.ForceColors)

	if *network {
		// create tendermint node
		node, err := tendermint.NewNode(cfg.Tendermint.Config, &cfg.Cosmos, dep.sdk.Service)
		if err != nil {
			logrus.WithField("module", "main").Fatalln(err)
		}
		logrus.WithField("module", "main").WithField("seeds", cfg.Tendermint.P2P.Seeds).Info("starting tendermint node")
		if err := node.Start(); err != nil {
			logrus.WithField("module", "main").Fatalln(err)
		}
	}

	// init system services.
	// if err := deployCoreServices(dep.cfg, dep.sdk); err != nil {
	// logrus.WithField("module", "main").Fatalln(err)
	// }

	// init gRPC server.
	server := grpc.New(dep.sdk)

	logrus.WithField("module", "main").Infof("starting MESG Engine version %s", version.Version)

	go func() {
		if err := server.Serve(cfg.Server.Address); err != nil {
			logrus.WithField("module", "main").Fatalln(err)
		}
	}()

	logrus.WithField("module", "main").Info("starting workflow engine")
	s := scheduler.New(dep.sdk.Event, dep.sdk.Execution, dep.sdk.Workflow)
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
	if err := stopRunningServices(dep.sdk); err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}
	server.Close()
}
