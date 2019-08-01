package workflow

import (
	"fmt"

	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	eventsdk "github.com/mesg-foundation/engine/sdk/event"
	executionsdk "github.com/mesg-foundation/engine/sdk/execution"
	servicesdk "github.com/mesg-foundation/engine/sdk/service"
	"github.com/mesg-foundation/engine/service"
)

// Workflow exposes functions of the workflow
type Workflow struct {
	event       *eventsdk.Event
	eventStream *eventsdk.Listener

	execution       *executionsdk.Execution
	executionStream *executionsdk.Listener

	service *servicesdk.Service

	ErrC chan error
}

// New creates a new Workflow instance
func New(event *eventsdk.Event, execution *executionsdk.Execution, service *servicesdk.Service) *Workflow {
	return &Workflow{
		event:     event,
		execution: execution,
		service:   service,
		ErrC:      make(chan error),
	}
}

// Start the workflow engine
func (w *Workflow) Start() error {
	if w.eventStream != nil || w.executionStream != nil {
		return fmt.Errorf("workflow engine already running")
	}
	w.eventStream = w.event.GetStream(nil)
	w.executionStream = w.execution.GetStream(&executionsdk.Filter{
		Statuses: []execution.Status{execution.Completed},
	})
	for {
		select {
		case event := <-w.eventStream.C:
			go w.processTrigger(service.EVENT, event.InstanceHash, event.Key, event.Data, event.Hash, nil)
		case execution := <-w.executionStream.C:
			go w.processTrigger(service.RESULT, execution.InstanceHash, execution.TaskKey, execution.Outputs, nil, execution.Hash)
		}
	}
}

func (w *Workflow) processTrigger(trigger service.TriggerType, instanceHash hash.Hash, key string, data map[string]interface{}, eventHash hash.Hash, executionHash hash.Hash) {
	services, err := w.service.List()
	if err != nil {
		w.ErrC <- err
		return
	}
	for _, service := range services {
		for _, wf := range service.Workflows {
			if wf.Trigger.Match(trigger, instanceHash, key, data) {
				_, err := w.execution.Execute(wf.Hash, wf.Tasks[0].InstanceHash, eventHash, executionHash, wf.Tasks[0].TaskKey, data, []string{})
				if err != nil {
					w.ErrC <- err
					continue
				}
			}
		}
	}
}
