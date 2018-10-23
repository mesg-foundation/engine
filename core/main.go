package main

import (
	"log"
	"path/filepath"

	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/interface/grpc"
	"github.com/mesg-foundation/core/logger"
	"github.com/mesg-foundation/core/systemservices"
	"github.com/mesg-foundation/core/version"
	"github.com/mesg-foundation/core/x/xsignal"
	"github.com/sirupsen/logrus"
)

func initGRPCServer(c *config.Config) (*grpc.Server, error) {
	// init services db.
	db, err := database.NewServiceDB(filepath.Join(c.Core.Path, c.Core.Database.RelativePath))
	if err != nil {
		return nil, err
	}

	// init api.
	a, err := api.New(db)
	if err != nil {
		return nil, err
	}

	return grpc.New(c.Server.Address, a), nil
}

func main() {
	// init configs.
	c, err := config.Global()
	if err != nil {
		log.Fatal(err)
	}

	// init logger.
	logger.Init(c.Log.Format, c.Log.Level)

	// init gRPC server.
	server, err := initGRPCServer(c)
	if err != nil {
		logrus.Fatalln(err)
	}

	logrus.Println("Starting MESG Core", version.Version)

	go func() {
		if err := server.Serve(); err != nil {
			logrus.Fatalln(err)
		}
	}()

	// TODO: rm this goroutine.
	// we have it here temporarily to test systemservices pkg.
	go func() {
		logrus.Println("service resolver: testing...")

		systemServicesPath := filepath.Join(c.Core.Path, c.SystemServices.RelativePath)

		a, err := api.New(db)
		if err != nil {
			logrus.Fatalln(err)
		}
		s, err := systemservices.New(a, systemServicesPath)
		if err != nil {
			logrus.Fatalln(err)
		}

		ss, validationErr, err := a.DeployServiceFromURL("https://github.com/mesg-foundation/service-webhook")
		if err != nil {
			logrus.Fatalln(err)
		}
		if validationErr != nil {
			logrus.Fatalln(err)
		}

		peers := []string{"core:50052"}
		logrus.Println("service resolver: webhook service started on peer:", peers[0])

		if err := s.Resolver().AddPeers(peers); err != nil {
			logrus.Fatalln(err)
		}
		logrus.Println("service resolver: webhook service's peer registered to resolver as:", peers[0])

		address, err := s.Resolver().Resolve(ss.ID)
		if err != nil {
			logrus.Fatalln(err)
		}
		logrus.Println("service resolver: peer address resolved for webhook service as:", address)
		logrus.Println("service resolver: works great!")
	}()

	<-xsignal.WaitForInterrupt()
	server.Close()
}
