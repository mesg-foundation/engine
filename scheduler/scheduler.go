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
			go s.processTrigger(workflow.RESULT, execution.InstanceHash, execution.TaskKey, execution.Outputs, nil, execution.Hash)
			go s.processExecution(execution)
		}
	}
}

func (s *Scheduler) processTrigger(trigger workflow.TriggerType, instanceHash hash.Hash, key string, data map[string]interface{}, eventHash hash.Hash, execHash hash.Hash) {
	workflows, err := s.workflow.List()
	if err != nil {
		s.ErrC <- err
		return
	}
	for _, wf := range workflows {
		if wf.Trigger.Match(trigger, instanceHash, key, data) {
			nextStep, err := wf.FindNode(wf.Trigger.InitialNode)
			if err != nil {
				s.ErrC <- err
				continue
			}
			if _, err := s.execution.Execute(wf.Hash, nextStep.InstanceHash, eventHash, execHash, wf.Trigger.InitialNode, nextStep.TaskKey, data, []string{}); err != nil {
				s.ErrC <- err
			}
		}
	}
}

func (s *Scheduler) processExecution(exec *execution.Execution) {
	if exec.WorkflowHash.IsZero() {
		return
	}
	wf, err := s.workflow.Get(exec.WorkflowHash)
	if err != nil {
		s.ErrC <- err
		return
	}
	for _, nextStepID := range wf.ChildrenIDs(exec.StepID) {
		nextStep, err := wf.FindNode(nextStepID)
		if err != nil {
			s.ErrC <- err
			continue
		}
		if _, err := s.execution.Execute(wf.Hash, nextStep.InstanceHash, nil, exec.Hash, nextStepID, nextStep.TaskKey, exec.Outputs, []string{}); err != nil {
			s.ErrC <- err
		}
	}
}
