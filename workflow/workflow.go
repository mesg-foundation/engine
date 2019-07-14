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
			go w.processEvent(event, Event)
		case execution := <-w.executionStream.C:
			eventHash, err := hash.Random()
			if err != nil {
				w.ErrC <- err
			}
			event := &event.Event{
				Hash:         eventHash,
				InstanceHash: execution.InstanceHash,
				Key:          execution.TaskKey,
				Data:         execution.Outputs,
			}
			go w.processEvent(event, Result)
		}
	}
}

func (w *Workflow) processEvent(event *event.Event, trigger triggerType) {
	workflows, err := w.findMatchingWorkflows(event, trigger)
	if err != nil {
		w.ErrC <- err
	}
	for _, workflow := range workflows {
		_, err := w.execution.Execute(workflow.Task.InstanceHash, event, workflow.Task.TaskKey, []string{})
		if err != nil {
			w.ErrC <- err
			continue
		}
	}
}

func (w *Workflow) findMatchingWorkflows(event *event.Event, trigger triggerType) ([]*workflow, error) {
	all, err := all()
	if err != nil {
		return nil, err
	}
	workflows := make([]*workflow, 0)
	for _, wf := range all {
		if wf.Trigger.MatchEvent(event) {
			workflows = append(workflows, wf)
		}
	}
	return workflows, nil
}

// All returns a fake set of data
// This is what can be called a system workflow and need to be removed when moved to services
// The hash of this instance correspond to the following service
// {"sid":"test-workflow","name":"Test workflow","tasks":[{"key":"taskX","inputs":[{"key":"foo","type":"String","object":[]},{"key":"bar","type":"String","object":[]}],"outputs":[{"key":"res","type":"Any","object":[]}]}],"events":[{"key":"eventX","data":[{"key":"foo","type":"String","object":[]},{"key":"bar","type":"String","object":[]}]}],"dependencies":[],"source":"QmQvRzJPFDhyBGK2rQP5mAeMrgp1XsTB8WYK2c7FHvyAB8"}
func all() ([]*workflow, error) {
	workflows := make([]*workflow, 0)
	instanceHash, err := hash.Decode("4fJs16kSV23Sc8CZ4nEJKoaQj1FogqWGrU2vpXT6vcbD")
	if err != nil {
		return nil, err
	}
	workflows = append(workflows, &workflow{
		Trigger: trigger{
			Key:          "taskX",
			Type:         Result,
			InstanceHash: instanceHash,
			Filters: []*filter{
				{Key: "bar", Predicate: EQ, Value: "world-2"},
			},
		},
		Task: task{
			InstanceHash: instanceHash,
			TaskKey:      "taskX",
		},
	})
	return workflows, nil
}
