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
	format, err := config.LogFormat().GetValue()
	if err != nil {
		panic(err)
	}
	level, err := config.LogLevel().GetValue()
	if err != nil {
		panic(err)
	}
	apiPort, err := config.APIPort().GetValue()
	if err != nil {
		panic(err)
	}

	logger.Init(format, level)

	logrus.Println("Starting MESG Core", version.Version)

	tcpServer := &grpc.Server{
		Network: "tcp",
		Address: ":" + apiPort,
	}

	go func() {
		if err := tcpServer.Serve(); err != nil {
			logrus.Fatalln(err)
		}
	}()

	<-xsignal.WaitForInterrupt()
	tcpServer.Close()
}
