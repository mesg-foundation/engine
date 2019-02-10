package main

import (
	"log"
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

func initGRPCServer(cfg *config.Config) (*grpc.Server, error) {
	// init services db.
	db, err := database.NewServiceDB(filepath.Join(cfg.Core.Path, cfg.Core.Database.ServiceRelativePath))
	if err != nil {
		return nil, err
	}

	// init execution db.
	execDB, err := database.NewExecutionDB(filepath.Join(cfg.Core.Path, cfg.Core.Database.ExecutionRelativePath))
	if err != nil {
		return nil, err
	}

	c, err := container.New()
	if err != nil {
		return nil, err
	}

	sm := service.NewContainerManager(c, cfg)

	// init api.
	api := api.New(db, execDB, sm)

	// init system services.
	if err := deployCoreServices(cfg, api); err != nil {
		return nil, err
	}

	return grpc.New(cfg.Server.Address, api), nil
}

func deployCoreServices(cfg *config.Config, api *api.API) error {
	for _, service := range cfg.Services() {
		logrus.Infof("Deploying service %q from %q", service.Key, service.URL)
		hash, err := api.DeployServiceFromURL(service.URL, service.Env, nil)
		if err != nil {
			return err
		}
		service.Hash = hash

		logrus.Infof("Service %q deployed", service.Key)
		if err := api.StartService(hash); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	// init configs.
	c, err := config.Global()
	if err != nil {
		log.Fatal(err)
	}

	// init logger.
	logger.Init(c.Log.Format, c.Log.Level, c.Log.ForceColors)

	// init gRPC server.
	server, err := initGRPCServer(c)
	if err != nil {
		logrus.Fatalln(err)
	}

	logrus.Infof("starting MESG Core version %s", version.Version)

	go func() {
		if err := server.Serve(); err != nil {
			logrus.Fatalln(err)
		}
	}()

	<-xsignal.WaitForInterrupt()
	server.Close()
}
