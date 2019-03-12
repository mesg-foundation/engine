package main

import (
	"path/filepath"

	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/interface/grpc"
	"github.com/mesg-foundation/core/logger"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/version"
	"github.com/mesg-foundation/core/x/xsignal"
	"github.com/sirupsen/logrus"
)

type dependencies struct {
	cfg         *config.Config
	serviceDB   database.ServiceDB
	executionDB database.ExecutionDB
	api         *api.API
}

func initDependencies() (*dependencies, error) {
	// init configs.
	cfg, err := config.Global()
	if err != nil {
		return nil, err
	}

	// init services db.
	serviceDB, err := database.NewServiceDB(filepath.Join(cfg.Core.Path, cfg.Core.Database.ServiceRelativePath))
	if err != nil {
		return nil, err
	}

	// init execution db.
	executionDB, err := database.NewExecutionDB(filepath.Join(cfg.Core.Path, cfg.Core.Database.ExecutionRelativePath))
	if err != nil {
		return nil, err
	}

	c, err := container.New()
	if err != nil {
		return nil, err
	}

	sm := service.NewContainerManager(c, cfg)

	// init api.
	api := api.New(serviceDB, executionDB, sm)

	return &dependencies{
		cfg:         cfg,
		serviceDB:   serviceDB,
		executionDB: executionDB,
		api:         api,
	}, nil
}

func deployCoreServices(cfg *config.Config, api *api.API) error {
	for _, service := range cfg.Services() {
		logrus.Infof("Deploying service %q from %q", service.Key, service.URL)
		s, err := api.DeployServiceFromURL(service.URL, service.Env, nil)
		if err != nil {
			return err
		}
		logrus.Infof("Service %q deployed", service.Key)
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
	for _, s := range services {
		if err := api.StopService(s.Hash); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	dep, err := initDependencies()
	if err != nil {
		logrus.Fatalln(err)
	}

	// init logger.
	logger.Init(dep.cfg.Log.Format, dep.cfg.Log.Level, dep.cfg.Log.ForceColors)

	// init system services.
	if err := deployCoreServices(dep.cfg, dep.api); err != nil {
		logrus.Fatalln(err)
	}

	// init gRPC server.
	server := grpc.New(dep.cfg.Server.Address, dep.api)

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
