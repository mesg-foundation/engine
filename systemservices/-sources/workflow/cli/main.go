package main

import (
	"log"

	"github.com/mesg-foundation/core/systemservices/-sources/workflow/workflow"
	"github.com/mesg-foundation/core/x/xsignal"
)

func main() {
	storage, err := workflow.NewMongoStorage("mongodb://mongo:27017", "workflow")
	if err != nil {
		log.Fatal(err)
	}

	// init WSS.
	w, err := workflow.New(storage)
	if err != nil {
		log.Fatal(err)
	}

	// start WSS.
	go func() {
		log.Println("WSS started")
		if err := w.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	// wait for interrupt and gracefully shutdown WSS.
	<-xsignal.WaitForInterrupt()
	log.Println("shutting down...")
	if err := w.Close(); err != nil {
		log.Fatal(err)
	}
	log.Println("shutdown")
}
