package main

import (
	"log"
	"path/filepath"

	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/interface/grpc"
	"github.com/mesg-foundation/core/logger"
	"github.com/mesg-foundation/core/version"
	"github.com/mesg-foundation/core/x/xsignal"
	"github.com/sirupsen/logrus"
)

func initGRPCServer(c *config.Config) (*grpc.Server, error) {
	// init services db.
	db, err := database.NewServiceDB(filepath.Join(c.Core.Path, c.Core.Database.ServiceRelativePath))
	if err != nil {
		return nil, err
	}

	// init execution db.
	execDB, err := database.NewExecutionDB(filepath.Join(c.Core.Path, c.Core.Database.ExecutionRelativePath))
	if err != nil {
		return nil, err
	}

	// init api.
	a, err := api.New(db, execDB)
	if err != nil {
		return nil, err
	}

	// init system services.
	if err := deployCoreServices(c, a); err != nil {
		return nil, err
	}

	return grpc.New(c.Server.Address, a), nil
}

func deployCoreServices(c *config.Config, api *api.API) error {
	for _, service := range c.Services() {
		logrus.Infof("Deploy service from %s", service.URL)
		s, valid, err := api.DeployServiceFromURL(service.URL, service.Env)
		if valid != nil {
			return valid
		}
		if err != nil {
			return err
		}
		logrus.Infof("Service %s deployed", s.Name)
		if err := api.StartService(s.Sid); err != nil {
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
