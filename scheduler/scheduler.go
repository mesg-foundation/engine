package scheduler

import (
	"fmt"

	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	eventsdk "github.com/mesg-foundation/engine/sdk/event"
	executionsdk "github.com/mesg-foundation/engine/sdk/execution"
	workflowsdk "github.com/mesg-foundation/engine/sdk/workflow"
	"github.com/mesg-foundation/engine/workflow"
)

// Scheduler manages the executions based on the definition of the workflows
type Scheduler struct {
	event       *eventsdk.Event
	eventStream *eventsdk.Listener

	execution       *executionsdk.Execution
	executionStream *executionsdk.Listener

	workflow *workflowsdk.Workflow

	ErrC chan error
}

// New creates a new Workflow instance
func New(event *eventsdk.Event, execution *executionsdk.Execution, workflow *workflowsdk.Workflow) *Scheduler {
	return &Scheduler{
		event:     event,
		execution: execution,
		workflow:  workflow,
		ErrC:      make(chan error),
	}
}

// Start the workflow engine
func (s *Scheduler) Start() error {
	if s.eventStream != nil || s.executionStream != nil {
		return fmt.Errorf("workflow scheduler already running")
	}
	s.eventStream = s.event.GetStream(nil)
	s.executionStream = s.execution.GetStream(&executionsdk.Filter{
		Statuses: []execution.Status{execution.Completed},
	})
	for {
		select {
		case event := <-s.eventStream.C:
			go s.processTrigger(workflow.EVENT, event.InstanceHash, event.Key, event.Data, event.Hash, nil)
		case execution := <-s.executionStream.C:
			go s.processTrigger(workflow.RESULT, execution.InstanceHash, execution.TaskKey, execution.Outputs, nil, execution)
			go s.processExecution(execution)
		}
	}
}

func (s *Scheduler) processTrigger(trigger workflow.TriggerType, instanceHash hash.Hash, key string, data map[string]interface{}, eventHash hash.Hash, exec *execution.Execution) {
	workflows, err := s.workflow.List()
	if err != nil {
		s.ErrC <- err
		return
	}
	for _, wf := range workflows {
		if wf.Trigger.Match(trigger, instanceHash, key, data) {
			if err := s.triggerExecution(wf, exec, eventHash, data); err != nil {
				s.ErrC <- err
			}
		}
	}
}

func (s *Scheduler) processExecution(exec *execution.Execution) error {
	if exec.WorkflowHash.IsZero() {
		return nil
	}
	wf, err := s.workflow.Get(exec.WorkflowHash)
	if err != nil {
		return err
	}
	return s.triggerExecution(wf, exec, nil, exec.Outputs)
}

func (s *Scheduler) triggerExecution(wf *workflow.Workflow, prev *execution.Execution, eventHash hash.Hash, data map[string]interface{}) error {
	height, err := s.getHeight(wf, prev)
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
	if _, err := s.execution.Execute(wf.Hash, task.InstanceHash, eventHash, parentHash, task.TaskKey, data, []string{}); err != nil {
		return err
	}
	return nil
}

func (s *Scheduler) getHeight(wf *workflow.Workflow, exec *execution.Execution) (int, error) {
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
	parent, err := s.execution.Get(exec.ParentHash)
	if err != nil {
		return 0, err
	}
	parentHeight, err := s.getHeight(wf, parent)
	return parentHeight + 1, err
}
