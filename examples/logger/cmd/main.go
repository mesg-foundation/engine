package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	mesg "github.com/mesg-foundation/go-service"
	"github.com/mesg-foundation/go-service/examples/logger/logger"
)

func main() {
	s, err := mesg.New()
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
