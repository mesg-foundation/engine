package main

import (
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/interface/grpc"
	"github.com/mesg-foundation/core/logger"
	"github.com/mesg-foundation/core/version"
	"github.com/mesg-foundation/core/x/xsignal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	format := viper.GetString(config.LogFormat)
	level := viper.GetString(config.LogLevel)
	logger.Init(format, level)

	logrus.Println("Starting MESG Core", version.Version)

	tcpServer := &grpc.Server{
		Network: "tcp",
		Address: viper.GetString(config.APIServerAddress),
	}

	unixServer := &grpc.Server{
		Network: "unix",
		Address: viper.GetString(config.APIServerSocket),
	}

	go startServer(tcpServer)
	go startServer(unixServer)

	closing := make(chan struct{}, 2)

	<-xsignal.WaitForInterrupt()

	go closeServer(tcpServer, closing)
	go closeServer(unixServer, closing)
	<-closing
	<-closing
}

func startServer(server *grpc.Server) {
	if err := server.Serve(); err != nil {
		logrus.Fatalln(err)
	}
}

func closeServer(server *grpc.Server, closing chan struct{}) {
	server.Close()
	closing <- struct{}{}
}
