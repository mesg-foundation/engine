package client

import (
	"context"
	"errors"
	"strings"

	"github.com/mesg-foundation/core/api/core"
)

// Start starts the workflow.
func (wf *Workflow) Start() error {
	if wf.Execute == nil {
		return errors.New("A workflow needs a task")
	}
	if wf.OnEvent == nil && wf.OnResult == nil {
		return errors.New("A workflow needs an event OnEvent or OnResult")
	}

	client, err := getClient()
	if err != nil {
		return err
	}
	wf.client = client

	if err := startServices(wf); err != nil {
		return err
	}

	listen := listenResults
	if wf.OnEvent != nil {
		listen = listenEvents
	}
	return listen(wf)
}

// Stop stops all the services in the workflow.
func (wf *Workflow) Stop() error {
	return stopServices(wf)
}

func listenEvents(wf *Workflow) error {
	if wf.OnEvent.Name == "" {
		return errors.New("Event's Name should be defined (you can use * to react to any event)")
	}
	stream, err := wf.client.ListenEvent(context.Background(), &core.ListenEventRequest{
		ServiceID: wf.OnEvent.ServiceID,
	})
	if err != nil {
		return err
	}

	for {
		data, err := stream.Recv()
		if err != nil {
			return err
		}
		if wf.validEvent(data) {
			if err := wf.Execute.processEvent(wf, data); err != nil {
				return err
			}
		}
	}
}

func (wf *Workflow) validEvent(data *core.EventData) bool {
	if strings.Compare(wf.OnEvent.Name, "*") == 0 {
		return true
	}
	return strings.Compare(wf.OnEvent.Name, data.EventKey) == 0
}

func listenResults(wf *Workflow) error {
	if wf.OnResult.Name == "" || wf.OnResult.Output == "" {
		return errors.New("Result's Name and Output should be defined (you can use * to react to any result)")
	}
	stream, err := wf.client.ListenResult(context.Background(), &core.ListenResultRequest{
		ServiceID: wf.OnResult.ServiceID,
	})
	if err != nil {
		return err
	}

	for {
		data, err := stream.Recv()
		if err != nil {
			return err
		}
		if wf.validResult(data) {
			if err := wf.Execute.processResult(wf, data); err != nil {
				return err
			}
		}
	}
}

func (wf *Workflow) validResult(data *core.ResultData) bool {
	return (wf.OnResult.Name == "*" || wf.OnResult.Name == data.TaskKey) &&
		(wf.OnResult.Output == "*" || wf.OnResult.Output == data.OutputKey)
}
