package main

import (
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
	"github.com/mesg-foundation/engine/tendermint/app"
	"github.com/mesg-foundation/engine/version"
	"github.com/mesg-foundation/engine/x/xerrors"
	"github.com/mesg-foundation/engine/x/xnet"
	"github.com/mesg-foundation/engine/x/xos"
	"github.com/mesg-foundation/engine/x/xsignal"
	"github.com/sirupsen/logrus"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/db"
)

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
	workflowDB, err := database.NewWorkflowDB(filepath.Join(config.Path, config.Database.WorkflowRelativePath))
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
		logrus.Infof("Deploying service %q", serviceConfig.Definition.Sid)
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
		logrus.Infof("Service %q deployed with hash %q", srv.Sid, srv.Hash)
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
		logrus.Infof("Instance started with hash %q", instance.Hash)
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
	cfg, err := config.Global()
	if err != nil {
		logrus.Fatalln(err)
	}

	dep, err := initDependencies(cfg)
	if err != nil {
		logrus.Fatalln(err)
	}

	// init logger.
	logger.Init(cfg.Log.Format, cfg.Log.Level, cfg.Log.ForceColors)

	// init app
	db := db.NewMemDB()
	logger := logger.TendermintLogger()
	app, _ := app.New(logger, db)

	// create tendermint node
	node, err := tendermint.NewNode(
		logger,
		app,
		filepath.Join(cfg.Path, cfg.Tendermint.Path),
		cfg.Tendermint.P2P.Seeds,
		cfg.Tendermint.P2P.ExternalAddress,
		ed25519.PubKeyEd25519(cfg.Tendermint.ValidatorPubKey),
	)
	if err != nil {
		logrus.Fatalln(err)
	}

	// init system services.
	if err := deployCoreServices(dep.cfg, dep.sdk); err != nil {
		logrus.Fatalln(err)
	}

	// init gRPC server.
	server := grpc.New(dep.sdk)

	logrus.Infof("starting MESG Engine version %s", version.Version)

	go func() {
		if err := server.Serve(cfg.Server.Address); err != nil {
			logrus.Fatalln(err)
		}
	}()

	logrus.WithField("seeds", cfg.Tendermint.P2P.Seeds).Info("starting tendermint node")
	go func() {
		if err := node.Start(); err != nil {
			logrus.Fatalln(err)
		}
		select {}
	}()

	logrus.Info("starting workflow engine")
	s := scheduler.New(dep.sdk.Event, dep.sdk.Execution, dep.sdk.Workflow)
	go func() {
		if err := s.Start(); err != nil {
			logrus.Fatalln(err)
		}
	}()
	go func() {
		for err := range s.ErrC {
			logrus.Warn(err)
		}
	}()

	<-xsignal.WaitForInterrupt()
	if err := stopRunningServices(dep.sdk); err != nil {
		logrus.Fatalln(err)
	}
	server.Close()
}
