package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mesg-foundation/core/client/service"
	"github.com/mesg-foundation/core/client/service/examples/logger/logger"
)

func main() {
	s, err := service.New()
	if err != nil {
		log.Fatal(err)
	}

	l := logger.New(s)
	if err := l.Start(); err != nil {
		log.Fatal(err)
	}

	abort := make(chan os.Signal, 1)
	signal.Notify(abort, syscall.SIGINT, syscall.SIGTERM)
	<-abort

	if err := l.Close(); err != nil {
		log.Fatal(err)
	}
}
