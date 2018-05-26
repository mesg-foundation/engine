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
	err := wf.startServices()
	if err != nil {
		panic(err)
	}

	stream, err := getClient().ListenEvent(context.Background(), &core.ListenEventRequest{
		ServiceID: wf.Event.Service,
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
		if wf.validEvent(data) {
			wf.processTasks(getClient(), data)
		}
	}
}

func (wf *Workflow) startServices() (err error) {
	return wf.iterateService(func(service string) (interface{}, error) {
		return getClient().StartService(context.Background(), &core.StartServiceRequest{
			ServiceID: service,
		})
	})
}

// Stop will stop all the services in your workflow
func (wf *Workflow) Stop() {
	wf.iterateService(func(service string) (interface{}, error) {
		return getClient().StopService(context.Background(), &core.StopServiceRequest{
			ServiceID: service,
		})
	})
}

func (wf *Workflow) validEvent(data *core.EventData) bool {
	return strings.Compare(data.EventKey, wf.Event.Name) == 0
}

func (wf *Workflow) processTasks(client core.CoreClient, data *core.EventData) (err error) {
	return wf.iterateTask(func(task *Task) error {
		return task.processEvent(client, data)
	})
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

func (wf *Workflow) iterateTask(process func(task *Task) error) (err error) {
	for _, task := range wf.Tasks {
		err = process(task)
		if err != nil {
			break
		}
	}
	return
}

func (wf *Workflow) iterateService(process func(service string) (interface{}, error)) (err error) {
	for _, service := range wf.getServices() {
		_, err = process(service)
		if err != nil {
			break
		}
	}
	return
}
