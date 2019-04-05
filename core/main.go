package main

import (
	"path/filepath"
	"sync"

	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/interface/grpc"
	"github.com/mesg-foundation/core/logger"
	"github.com/mesg-foundation/core/version"
	"github.com/mesg-foundation/core/x/xerrors"
	"github.com/mesg-foundation/core/x/xsignal"
	"github.com/sirupsen/logrus"
)

type dependencies struct {
	config      *config.Config
	serviceDB   database.ServiceDB
	executionDB database.ExecutionDB
	container   container.Container
	api         *api.API
}

func initDependencies() (*dependencies, error) {
	// init configs.
	config, err := config.Global()
	if err != nil {
		return nil, err
	}

	// init services db.
	serviceDB, err := database.NewServiceDB(filepath.Join(config.Core.Path, config.Core.Database.ServiceRelativePath))
	if err != nil {
		return nil, err
	}

	// init execution db.
	executionDB, err := database.NewExecutionDB(filepath.Join(config.Core.Path, config.Core.Database.ExecutionRelativePath))
	if err != nil {
		return nil, err
	}

	// init container.
	c, err := container.New()
	if err != nil {
		return nil, err
	}

	// init api.
	api := api.New(c, serviceDB, executionDB)

	return &dependencies{
		config:      config,
		serviceDB:   serviceDB,
		executionDB: executionDB,
		container:   c,
		api:         api,
	}, nil
}

func deployCoreServices(config *config.Config, api *api.API) error {
	for _, service := range config.Services() {
		logrus.Infof("Deploying service %q from %q", service.Key, service.URL)
		s, valid, err := api.DeployServiceFromURL(service.URL, service.Env)
		if valid != nil {
			return valid
		}
		if err != nil {
			return err
		}
		service.Sid = s.Sid
		service.Hash = s.Hash
		logrus.Infof("Service %q deployed with hash %q", service.Key, service.Hash)
		if err := api.StartService(s.Sid); err != nil {
			return err
		}
	}
	return nil
}

func stopRunningServices(api *api.API) error {
	services, err := api.ListServices()
	if err != nil {
		return err
	}
	var (
		serviceLen = len(services)
		errC       = make(chan error, serviceLen)
		wg         sync.WaitGroup
	)
	wg.Add(serviceLen)
	for _, service := range services {
		go func(hash string) {
			defer wg.Done()
			err := api.StopService(hash)
			if err != nil {
				errC <- err
			}
		}(service.Hash)
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
	if err := deployCoreServices(dep.config, dep.api); err != nil {
		logrus.Fatalln(err)
	}

	// init gRPC server.
	server := grpc.New(dep.config.Server.Address, dep.container, dep.api)

	logrus.Infof("starting MESG Core version %s", version.Version)

	go func() {
		if err := server.Serve(); err != nil {
			logrus.Fatalln(err)
		}
	}()

	<-xsignal.WaitForInterrupt()
	if err := stopRunningServices(dep.api); err != nil {
		logrus.Fatalln(err)
	}
	server.Close()
}
