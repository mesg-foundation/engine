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
	if wf.OnEvent == nil {
		panic("A workflow needs an event")
	}
	startServices(wf)
	listenEvents(wf)
}

// Stop will stop all the services in your workflow
func (wf *Workflow) Stop() {
	stopServices(wf)
}

func listenEvents(wf *Workflow) {
	stream, err := getClient().ListenEvent(context.Background(), &core.ListenEventRequest{
		ServiceID: wf.OnEvent.Service,
	})
	if err != nil {
		panic(err)
	}

	for {
		var data *core.EventData
		data, err = stream.Recv()
		if err != nil {
			panic(err)
		}
		log.Println("Event received", data)
		if strings.Compare(data.EventKey, wf.OnEvent.Name) == 0 {
			wf.Execute.processEvent(data)
		}
	}
}
