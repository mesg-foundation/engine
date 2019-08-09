package engine

import (
	"fmt"

	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	eventsdk "github.com/mesg-foundation/engine/sdk/event"
	executionsdk "github.com/mesg-foundation/engine/sdk/execution"
	workflowsdk "github.com/mesg-foundation/engine/sdk/workflow"
	"github.com/mesg-foundation/engine/workflow"
)

// Workflow exposes functions of the workflow
type Workflow struct {
	event       *eventsdk.Event
	eventStream *eventsdk.Listener

	execution       *executionsdk.Execution
	executionStream *executionsdk.Listener

	workflow *workflowsdk.Workflow

	ErrC chan error
}

// New creates a new Workflow instance
func New(event *eventsdk.Event, execution *executionsdk.Execution, workflow *workflowsdk.Workflow) *Workflow {
	return &Workflow{
		event:     event,
		execution: execution,
		workflow:  workflow,
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
			go w.processTrigger(workflow.EVENT, event.InstanceHash, event.Key, event.Data, event.Hash, nil)
		case execution := <-w.executionStream.C:
			go w.processTrigger(workflow.RESULT, execution.InstanceHash, execution.TaskKey, execution.Outputs, nil, execution)
			go w.processExecution(execution)
		}
	}
}

func (w *Workflow) processTrigger(trigger workflow.TriggerType, instanceHash hash.Hash, key string, data map[string]interface{}, eventHash hash.Hash, exec *execution.Execution) {
	workflows, err := w.workflow.List()
	if err != nil {
		w.ErrC <- err
		return
	}
	for _, wf := range workflows {
		if wf.Trigger.Match(trigger, instanceHash, key, data) {
			if err := w.triggerExecution(wf, exec, eventHash, data); err != nil {
				w.ErrC <- err
			}
		}
	}
}

func (w *Workflow) processExecution(exec *execution.Execution) error {
	if exec.WorkflowHash.IsZero() {
		return nil
	}
	wf, err := w.workflow.Get(exec.WorkflowHash)
	if err != nil {
		return err
	}
	return w.triggerExecution(wf, exec, nil, exec.Outputs)
}

func (w *Workflow) triggerExecution(wf *workflow.Workflow, prev *execution.Execution, eventHash hash.Hash, data map[string]interface{}) error {
	height, err := w.getHeight(wf, prev)
	if err != nil {
		return err
	}
	if len(wf.Tasks) <= height {
		// end of workflow
		return nil
	}
	var parentHash hash.Hash
	if prev != nil {
		parentHash = prev.Hash
	}
	task := wf.Tasks[height]
	if _, err := w.execution.Execute(wf.Hash, task.InstanceHash, eventHash, parentHash, task.TaskKey, data, []string{}); err != nil {
		return err
	}
	return nil
}

func (w *Workflow) getHeight(wf *workflow.Workflow, exec *execution.Execution) (int, error) {
	if exec == nil {
		return 0, nil
	}
	// Result from other workflow
	if !exec.WorkflowHash.Equal(wf.Hash) {
		return 0, nil
	}
	// Execution triggered by an event
	if !exec.EventHash.IsZero() {
		return 1, nil
	}
	if exec.ParentHash.IsZero() {
		panic("parent hash should be present if event is not")
	}
	if exec.ParentHash.Equal(exec.Hash) {
		panic("parent hash cannot be equal to execution hash")
	}
	parent, err := w.execution.Get(exec.ParentHash)
	if err != nil {
		return 0, err
	}
	parentHeight, err := w.getHeight(wf, parent)
	return parentHeight + 1, err
}
