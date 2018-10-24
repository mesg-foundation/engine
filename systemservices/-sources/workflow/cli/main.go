package main

import (
	"log"

	"github.com/mesg-foundation/core/systemservices/-sources/workflow/workflow"
	"github.com/mesg-foundation/core/x/xsignal"
	"github.com/sirupsen/logrus"
)

var (
	coreAddr  = "core:50052"
	mongoAddr = "mongodb://mongo:27017"
	mongoDB   = "workflow"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	storage, err := workflow.NewMongoStorage(mongoAddr, mongoDB)
	if err != nil {
		log.Fatal(err)
	}

	// init WSS.
	w, err := workflow.New(coreAddr, storage)
	if err != nil {
		log.Fatal(err)
	}

	// start WSS.
	go func() {
		logrus.WithFields(logrus.Fields{
			"general": true,
		}).Info("WSS started")

		if err := w.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	// wait for interrupt and gracefully shutdown WSS.
	<-xsignal.WaitForInterrupt()

	logrus.WithFields(logrus.Fields{
		"general": true,
	}).Info("shutting down...")

	if err := w.Close(); err != nil {
		log.Fatal(err)
	}

	logrus.WithFields(logrus.Fields{
		"general": true,
	}).Info("shutdown")
}
