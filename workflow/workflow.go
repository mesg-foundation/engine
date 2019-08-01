package workflow

import (
	"fmt"

	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	eventsdk "github.com/mesg-foundation/engine/sdk/event"
	executionsdk "github.com/mesg-foundation/engine/sdk/execution"
	servicesdk "github.com/mesg-foundation/engine/sdk/service"
	"github.com/mesg-foundation/engine/service"
	"github.com/sirupsen/logrus"
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
			go w.processTrigger(service.RESULT, execution.InstanceHash, execution.TaskKey, execution.Outputs, nil, execution)
			go w.processExecution(execution)
		}
	}
}

func (w *Workflow) processTrigger(trigger service.TriggerType, instanceHash hash.Hash, key string, data map[string]interface{}, eventHash hash.Hash, exec *execution.Execution) {
	services, err := w.service.List()
	if err != nil {
		w.ErrC <- err
		return
	}
	for _, service := range services {
		for _, wf := range service.Workflows {
			if wf.Trigger.Match(trigger, instanceHash, key, data) {
				if err := w.triggerExecution(wf, exec, eventHash, data); err != nil {
					w.ErrC <- err
				}
			}
		}
	}
}

func (w *Workflow) processExecution(exec *execution.Execution) error {
	if exec.WorkflowHash.IsZero() {
		return nil
	}
	wf, err := w.service.FindWorkflow(exec.WorkflowHash)
	if err != nil {
		return err
	}
	return w.triggerExecution(wf, exec, nil, exec.Outputs)
}

func (w *Workflow) triggerExecution(wf *service.Workflow, prev *execution.Execution, eventHash hash.Hash, data map[string]interface{}) error {
	height, err := w.getHeight(prev)
	if err != nil {
		return err
	}
	if len(wf.Tasks) <= height {
		// end of workflow
		return nil
	}
	task := wf.Tasks[height]
	hash, err := w.execution.Execute(wf.Hash, task.InstanceHash, eventHash, prev.Hash, task.TaskKey, data, []string{})
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"workflow": wf.Key,
		"task":     task.TaskKey,
		"exec":     hash.String(),
		"parent":   prev.Hash.String(),
	}).Debug("workflow execution")
	return nil
}

func (w *Workflow) getHeight(exec *execution.Execution) (int, error) {
	if exec == nil {
		return 0, nil
	}
	if exec.ParentHash == nil {
		return 0, nil
	}
	parent, err := w.execution.Get(exec.ParentHash)
	if err != nil {
		return 0, err
	}
	parentHeight, err := w.getHeight(parent)
	return parentHeight + 1, err
}
