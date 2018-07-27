// Package main is a quick-start application. Please visit:
// https://docs.mesg.com/start-here/quick-start-guide
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ilgooz/mesg-go/application"
	"github.com/ilgooz/mesg-go/examples/application-quickstart/quickstart"
)

var config = quickstart.Config{
	WebhookServiceID:    "v1_61a1850a786dc8c2d59e8d7d5aaaecce",
	DiscordInvServiceID: "v1_a9ac300dcd1322205b6aa00a382f5775",
	LogServiceID:        "v1_ea41381d7dbd9ad6869827a89aad2f25",
	SendgridKey:         "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",
	Email:               "ilkergoktugozturk@gmail.com",
}

func main() {
	app, err := application.New()
	if err != nil {
		log.Fatal(err)
	}

	q := quickstart.New(app, config)
	go func() {
		if err := q.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	abort := make(chan os.Signal, 1)
	signal.Notify(abort, syscall.SIGINT, syscall.SIGTERM)
	<-abort

	if err := q.Close(); err != nil {
		log.Fatal(err)
	}
}
