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
		Address: ":" + viper.GetString(config.APIGRPCPort),
	}

	go func() {
		if err := tcpServer.Serve(); err != nil {
			logrus.Fatalln(err)
		}
	}()

	<-xsignal.WaitForInterrupt()
	tcpServer.Close()
}
