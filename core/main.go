package main

import (
	"path/filepath"
	"sync"

	"github.com/mesg-foundation/engine/config"
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/database"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/logger"
	"github.com/mesg-foundation/engine/sdk"
	instancesdk "github.com/mesg-foundation/engine/sdk/instance"
	servicesdk "github.com/mesg-foundation/engine/sdk/service"
	"github.com/mesg-foundation/engine/server/grpc"
	"github.com/mesg-foundation/engine/version"
	"github.com/mesg-foundation/engine/workflow"
	"github.com/mesg-foundation/engine/x/xerrors"
	"github.com/mesg-foundation/engine/x/xos"
	"github.com/mesg-foundation/engine/x/xsignal"
	"github.com/sirupsen/logrus"
)

type dependencies struct {
	config      *config.Config
	serviceDB   database.ServiceDB
	executionDB database.ExecutionDB
	container   container.Container
	sdk         *sdk.SDK
}

func initDependencies() (*dependencies, error) {
	// init configs.
	config, err := config.Global()
	if err != nil {
		return nil, err
	}

	// init services db.
	serviceDB, err := database.NewServiceDB(filepath.Join(config.Path, config.Database.ServiceRelativePath))
	if err != nil {
		return nil, err
	}

	// init instance db.
	instanceDB, err := database.NewInstanceDB(filepath.Join(config.Path, config.Database.InstanceRelativePath))
	if err != nil {
		return nil, err
	}

	// init execution db.
	executionDB, err := database.NewExecutionDB(filepath.Join(config.Path, config.Database.ExecutionRelativePath))
	if err != nil {
		return nil, err
	}

	// init container.
	c, err := container.New(config.Name)
	if err != nil {
		return nil, err
	}

	// init sdk.
	sdk := sdk.New(c, serviceDB, instanceDB, executionDB)

	return &dependencies{
		config:      config,
		container:   c,
		serviceDB:   serviceDB,
		executionDB: executionDB,
		sdk:         sdk,
	}, nil
}

func deployCoreServices(config *config.Config, sdk *sdk.SDK) error {
	for _, serviceConfig := range config.SystemServices {
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
	dep, err := initDependencies()
	if err != nil {
		logrus.Fatalln(err)
	}

	// init logger.
	logger.Init(dep.config.Log.Format, dep.config.Log.Level, dep.config.Log.ForceColors)

	// init system services.
	if err := deployCoreServices(dep.config, dep.sdk); err != nil {
		logrus.Fatalln(err)
	}

	// init gRPC server.
	server := grpc.New(dep.sdk)

	logrus.Infof("starting MESG Engine version %s", version.Version)

	go func() {
		if err := server.Serve(dep.config.Server.Address); err != nil {
			logrus.Fatalln(err)
		}
	}()

	go func() {
		logrus.Info("starting workflow engine")
		wf := workflow.New(dep.sdk.Event, dep.sdk.Execution)
		if err := wf.Start(); err != nil {
			logrus.Fatalln(err)
		}
	}()

	<-xsignal.WaitForInterrupt()
	if err := stopRunningServices(dep.sdk); err != nil {
		logrus.Fatalln(err)
	}
	server.Close()
}
