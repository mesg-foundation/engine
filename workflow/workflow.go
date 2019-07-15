package workflow

import (
	"fmt"

	"github.com/mesg-foundation/engine/event"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	eventsdk "github.com/mesg-foundation/engine/sdk/event"
	executionsdk "github.com/mesg-foundation/engine/sdk/execution"
)

// Workflow exposes functions of the workflow
type Workflow struct {
	event       *eventsdk.Event
	eventStream *eventsdk.Listener

	execution       *executionsdk.Execution
	executionStream *executionsdk.Listener

	ErrC chan error
}

// New creates a new Workflow instance
func New(event *eventsdk.Event, execution *executionsdk.Execution) *Workflow {
	return &Workflow{
		event:     event,
		execution: execution,
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
			go w.processEvent(event)
		case execution := <-w.executionStream.C:
			go w.processExecution(execution)
		}
	}
}

func (w *Workflow) processEvent(event *event.Event) {
	all, err := all()
	if err != nil {
		w.ErrC <- err
		return
	}
	for _, wf := range all {
		if wf.Trigger.Type == Event && wf.Trigger.InstanceHash.Equal(event.InstanceHash) && wf.Trigger.Key == event.Key {
			_, err := w.execution.Execute(wf.Task.InstanceHash, event.Hash, nil, wf.Task.TaskKey, event.Data, []string{})
			if err != nil {
				w.ErrC <- err
				continue
			}
		}
	}
}

func (w *Workflow) processExecution(execution *execution.Execution) {
	all, err := all()
	if err != nil {
		w.ErrC <- err
		return
	}
	for _, wf := range all {
		if wf.Trigger.MatchEvent(event) {
			_, err := w.execution.Execute(wf.Task.InstanceHash, nil, execution.Hash, wf.Task.TaskKey, execution.Outputs, []string{})
			if err != nil {
				w.ErrC <- err
				continue
			}
		}
	}
}

// All returns a fake set of data
// This is what can be called a system workflow and need to be removed when moved to services
// The hash of this instance correspond to the following service
// {"sid":"test-workflow","name":"test-workflow","tasks":[{"key":"taskA","inputs":[{"key":"a","type":"String","object":[]}],"outputs":[{"key":"a","type":"String","object":[]},{"key":"b","type":"Boolean","object":[]}]},{"key":"taskB","inputs":[{"key":"a","type":"String","object":[]},{"key":"b","type":"Boolean","object":[]}],"outputs":[{"key":"a","type":"Boolean","object":[]}]}],"events":[{"key":"started","data":[{"key":"a","type":"String","object":[]}]},{"key":"interval","data":[{"key":"i","type":"Number","object":[]}]}],"dependencies":[],"configuration":{"env":["EVENT_INTERVAL=1000"]},"source":"QmVvDmnTWnUd4EmWT3qUoBiv8gym4EYkTkwfBmfvYs7WFS"}
func all() ([]*workflow, error) {
	instanceHash, err := hash.Decode("FtxZoLSD4M8w4v3ZfY8s6tY4bAc9B3Wy1TiqGq8iP3Tt")
	if err != nil {
		return nil, err
	}
	return []*workflow{
		{ // When result of taskA() -> (a string, b string), execute taskB(a string, b bool)
			Trigger: trigger{InstanceHash: instanceHash, Type: Result, Key: "taskA"},
			Task:    task{InstanceHash: instanceHash, TaskKey: "taskB"},
		},
		{ // When event started(a string), execute taskA(a string)
			Trigger: trigger{InstanceHash: instanceHash, Type: Event, Key: "started"},
			Task:    task{InstanceHash: instanceHash, TaskKey: "taskA"},
		},
	}, nil
}
