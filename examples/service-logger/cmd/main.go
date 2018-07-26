package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ilgooz/mesg-go/service"
	"github.com/ilgooz/service-logger/logger"
)

func main() {
	s, err := service.New(
		service.EndpointOption("endpoint"),
		service.TokenOption("endpoint"),
	)
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
