package client

import (
	"context"
	"log"
	"strings"

	"github.com/mesg-foundation/core/api/core"
)

// Then connects a task to a workflow
func (wf *Workflow) Then(task *Task) *Workflow {
	wf.Tasks = append(wf.Tasks, task)
	return wf
}

// Start is the function to start the workflow
func (wf *Workflow) Start() {
	var err error
	client := getClient()
	for _, service := range wf.getServices() {
		_, err = client.StartService(context.Background(), &core.StartServiceRequest{
			ServiceID: service,
		})
		if err != nil {
			panic(err)
		}
	}

	stream, err := client.ListenEvent(context.Background(), &core.ListenEventRequest{
		ServiceID: wf.Event.Service,
	})
	if err != nil {
		panic(err)
	}

	for {
		var data *core.EventData
		data, err = stream.Recv()
		log.Println("Event received", data)
		if err != nil {
			panic(err)
		}
		if strings.Compare(data.EventKey, wf.Event.Name) == 0 {
			wf.processTasks(client, data)
		}
	}
}

// Stop will stop all the services in your workflow
func (wf *Workflow) Stop() {
	client := getClient()
	for _, service := range wf.getServices() {
		client.StopService(context.Background(), &core.StopServiceRequest{
			ServiceID: service,
		})
	}
}

func (wf *Workflow) processTasks(client core.CoreClient, data *core.EventData) (err error) {
	for _, task := range wf.Tasks {
		err = task.processEvent(client, data)
		if err != nil {
			break
		}
	}
	return
}

func (wf *Workflow) getServices() (services []string) {
	services = make([]string, 0)
	addedServices := make(map[string]bool)
	if !addedServices[wf.Event.Service] {
		addedServices[wf.Event.Service] = true
		services = append(services, wf.Event.Service)
	}
	for _, t := range wf.Tasks {
		if !addedServices[t.Service] {
			addedServices[t.Service] = true
			services = append(services, t.Service)
		}
	}
	return
}
