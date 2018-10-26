package main

import (
	"path/filepath"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/interface/grpc"
	"github.com/mesg-foundation/core/logger"
	"github.com/mesg-foundation/core/version"
	"github.com/mesg-foundation/core/x/xsignal"
	"github.com/sirupsen/logrus"
)

func main() {
	c, err := config.Global()
	if err != nil {
		panic(err)
	}

	db, err := database.NewServiceDB(filepath.Join(c.Core.Path, c.Core.Database.ServiceRelativePath))
	if err != nil {
		logrus.Fatalln(err)
	}

	execDB, err := execution.New(filepath.Join(c.Core.Path, c.Core.Database.ExecutionRelativePath))
	if err != nil {
		logrus.Fatalln(err)
	}

	logger.Init(c.Log.Format, c.Log.Level)

	logrus.Println("Starting MESG Core", version.Version)

	tcpServer := &grpc.Server{
		Network:   "tcp",
		Address:   c.Server.Address,
		ServiceDB: db,
		ExecDB:    execDB,
	}

	go func() {
		if err := tcpServer.Serve(); err != nil {
			logrus.Fatalln(err)
		}
	}()

	<-xsignal.WaitForInterrupt()
	tcpServer.Close()
}
