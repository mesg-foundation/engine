package main

import (
	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/config"
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
	go startServer(&api.Server{
		Network: "tcp",
		Address: viper.GetString(config.APIServerAddress),
	})
	go startServer(&api.Server{
		Network: "unix",
		Address: viper.GetString(config.APIServerSocket),
	})
	<-xsignal.WaitForInterrupt()
}

func startServer(server *api.Server) {
	err := server.Serve()
	defer server.Stop()
	if err != nil {
		logrus.Fatalln(err)
	}
}
