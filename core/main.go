package main

import (
	"path/filepath"

	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/interface/grpc"
	"github.com/mesg-foundation/core/logger"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/version"
	"github.com/mesg-foundation/core/x/xsignal"
	"github.com/sirupsen/logrus"
)

type dependencies struct {
	config      *config.Config
	serviceDB   database.ServiceDB
	executionDB database.ExecutionDB
	api         *api.API
}

func initDependencies() (*dependencies, error) {
	dep := dependencies{}

	// init configs.
	config, err := config.Global()
	if err != nil {
		return &dep, err
	}
	dep.config = config

	// init services db.
	serviceDB, err := database.NewServiceDB(filepath.Join(dep.config.Core.Path, dep.config.Core.Database.ServiceRelativePath))
	if err != nil {
		return &dep, err
	}
	dep.serviceDB = serviceDB

	// init execution db.
	executionDB, err := database.NewExecutionDB(filepath.Join(dep.config.Core.Path, dep.config.Core.Database.ExecutionRelativePath))
	if err != nil {
		return &dep, err
	}
	dep.executionDB = executionDB

	// init api.
	api, err := api.New(dep.serviceDB, dep.executionDB)
	if err != nil {
		return &dep, err
	}
	dep.api = api

	return &dep, nil
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
		status, err := s.Status()
		if err != nil {
			return err
		}
		if status == service.STARTING || status == service.PARTIAL || status == service.RUNNING {
			if err := api.StopService(s.Hash); err != nil {
				return err
			}
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
	logger.Init(dep.config.Log.Format, dep.config.Log.Level, dep.config.Log.ForceColors)

	// init system services.
	if err := deployCoreServices(dep.config, dep.api); err != nil {
		logrus.Fatalln(err)
	}

	// init gRPC server.
	server := grpc.New(dep.config.Server.Address, dep.api)

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
