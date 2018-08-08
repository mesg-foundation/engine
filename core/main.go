package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/version"
	"github.com/spf13/viper"
)

func main() {
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
		log.Panicln(err)
	}
}
