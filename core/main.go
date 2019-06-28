package main

import (
	"path/filepath"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/logger"
	"github.com/mesg-foundation/core/sdk"
	"github.com/mesg-foundation/core/server/grpc"
	"github.com/mesg-foundation/core/version"
	"github.com/mesg-foundation/core/x/xsignal"
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
	c, err := container.New()
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

// func deployCoreServices(config *config.Config, sdk *sdk.SDK) error {
// 	for _, service := range config.Services() {
// 		logrus.Infof("Deploying service %q from %q", service.Key, service.URL)
// 		s, valid, err := sdk.DeployServiceFromURL(service.URL, service.Env)
// 		if valid != nil {
// 			return valid
// 		}
// 		if err != nil {
// 			return err
// 		}
// 		service.Sid = s.Sid
// 		service.Hash = s.Hash
// 		logrus.Infof("Service %q deployed with hash %q", service.Key, service.Hash)
// 		if err := sdk.StartService(s.Hash); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func stopRunningServices(sdk *sdk.SDK) error {
// 	services, err := sdk.ListServices()
// 	if err != nil {
// 		return err
// 	}
// 	var (
// 		serviceLen = len(services)
// 		errC       = make(chan error, serviceLen)
// 		wg         sync.WaitGroup
// 	)
// 	wg.Add(serviceLen)
// 	for _, service := range services {
// 		go func(hash hash.Hash) {
// 			defer wg.Done()
// 			err := sdk.StopService(hash)
// 			if err != nil {
// 				errC <- err
// 			}
// 		}(service.Hash)
// 	}
// 	wg.Wait()
// 	close(errC)
// 	var errs xerrors.Errors
// 	for err := range errC {
// 		errs = append(errs, err)
// 	}
// 	return errs.ErrorOrNil()
// }

func main() {
	dep, err := initDependencies()
	if err != nil {
		logrus.Fatalln(err)
	}

	// init logger.
	logger.Init(dep.config.Log.Format, dep.config.Log.Level, dep.config.Log.ForceColors)

	// init system services.
	// if err := deployCoreServices(dep.config, dep.sdk); err != nil {
	// 	logrus.Fatalln(err)
	// }

	// init gRPC server.
	server := grpc.New(dep.sdk)

	logrus.Infof("starting MESG Engine version %s", version.Version)

	go func() {
		if err := server.Serve(dep.config.Server.Address); err != nil {
			logrus.Fatalln(err)
		}
	}()

	<-xsignal.WaitForInterrupt()
	// if err := stopRunningServices(dep.sdk); err != nil {
	// 	logrus.Fatalln(err)
	// }
	server.Close()
}
