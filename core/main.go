package main

import (
	"github.com/mesg-foundation/core/config"
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

	logger.Init(c.Log.Format, c.Log.Level)

	logrus.Println("Starting MESG Core", version.Version)

	tcpServer := &grpc.Server{
		Network: "tcp",
		Address: c.Server.Address,
	}

	go func() {
		if err := tcpServer.Serve(); err != nil {
			logrus.Fatalln(err)
		}
	}()

	<-xsignal.WaitForInterrupt()
	tcpServer.Close()
}
