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
	WebhookServiceID:    "v1_3b56fdfc5bb5f50def11efaa6045046c",
	DiscordInvServiceID: "v1_e4e8d488444bbd9389ad9eac150669be",
	LogServiceID:        "v1_0df94b3537841789f1dbe8bc8690a4c6",
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
