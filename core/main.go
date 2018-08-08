package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/logger"
	"github.com/mesg-foundation/core/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	format := viper.GetString(config.LogFormat)
	level := viper.GetString(config.LogLevel)
	logger.Init(format, level)

	log.Println("Starting MESG Core", version.Version)
	go startServer(&api.Server{
		Network: "tcp",
		Address: viper.GetString(config.APIServerAddress),
	})
	go startServer(&api.Server{
		Network: "unix",
		Address: viper.GetString(config.APIServerSocket),
	})
	abort := make(chan os.Signal, 1)
	signal.Notify(abort, syscall.SIGINT, syscall.SIGTERM)
	<-abort
}

func startServer(server *api.Server) {
	err := server.Serve()
	defer server.Stop()
	if err != nil {
		log.Fatalln(err)
	}
}
