package client

import (
	"context"
	"log"
	"strings"

	"github.com/mesg-foundation/core/api/core"
)

// Start is the function to start the workflow
func (wf *Workflow) Start() {
	if wf.Execute == nil {
		panic("A workflow needs a taks")
	}
	if wf.OnEvent == nil && wf.OnResult == nil {
		panic("A workflow needs an event OnEvent or OnResult")
	}
	startServices(wf)
	if wf.OnEvent != nil {
		listenEvents(wf)
	} else {
		listenResults(wf)
	}
}

// Stop will stop all the services in your workflow
func (wf *Workflow) Stop() {
	stopServices(wf)
}

func listenEvents(wf *Workflow) {
	if wf.OnEvent.Name == "" {
		panic("Event's Name should be defined (you can use * to react to any event)")
	}
	stream, err := getClient().ListenEvent(context.Background(), &core.ListenEventRequest{
		ServiceID: wf.OnEvent.Service,
	})
	if err != nil {
		panic(err)
	}
	log.Println("Listening events from", wf.OnEvent.Name, "...")

	for {
		var data *core.EventData
		data, err = stream.Recv()
		if err != nil {
			panic(err)
		}
		log.Println("Event received", data)
		if strings.Compare(data.EventKey, wf.OnEvent.Name) == 0 || wf.OnEvent.Name == "*" {
			wf.Execute.processEvent(data)
		}
	}
}

func listenResults(wf *Workflow) {
	if wf.OnResult.Name == "" || wf.OnResult.Output == "" {
		panic("Result's Name and Output should be defined (you can use * to react to any result)")
	}
	stream, err := getClient().ListenResult(context.Background(), &core.ListenResultRequest{
		ServiceID: wf.OnResult.Service,
	})
	if err != nil {
		panic(err)
	}
	log.Println("Listening results from", wf.OnResult.Name, wf.OnResult.Output, "...")

	for {
		var data *core.ResultData
		data, err = stream.Recv()
		if err != nil {
			panic(err)
		}
		log.Println("Result received", data)
		if (strings.Compare(data.TaskKey, wf.OnResult.Name) == 0 || wf.OnResult.Name == "*") &&
			(strings.Compare(data.OutputKey, wf.OnResult.Output) == 0 || wf.OnResult.Output == "*") {
			wf.Execute.processResult(data)
		}
	}
}
