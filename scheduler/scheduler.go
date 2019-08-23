package scheduler

import (
	"fmt"

	"github.com/mesg-foundation/engine/event"
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
			go s.processTriggerFromEvent(event)
		case execution := <-s.executionStream.C:
			go s.processTriggerFromResult(execution)
			go s.processExecution(execution)
		}
	}
}

func (s *Scheduler) processTriggerFromEvent(event *event.Event) {
	workflows, err := s.workflowsMatchingFilter(func(wf *workflow.Workflow) bool {
		return wf.Trigger.InstanceHash.Equal(event.InstanceHash) &&
			wf.Trigger.EventKey == event.Key &&
			wf.Trigger.Filters.Match(event.Data)
	})
	if err != nil {
		s.ErrC <- err
		return
	}
	for _, wf := range workflows {
		nextStep, err := wf.FindNode(wf.Trigger.NodeKey)
		if err != nil {
			s.ErrC <- err
			continue
		}
		if _, err := s.execution.Execute(wf.Hash, nextStep.InstanceHash, event.Hash, nil, wf.Trigger.NodeKey, nextStep.TaskKey, event.Data, []string{}); err != nil {
			s.ErrC <- err
		}
	}
}

func (s *Scheduler) processTriggerFromResult(result *execution.Execution) {
	workflows, err := s.workflowsMatchingFilter(func(wf *workflow.Workflow) bool {
		return wf.Trigger.InstanceHash.Equal(result.InstanceHash) &&
			wf.Trigger.TaskKey == result.TaskKey &&
			wf.Trigger.Filters.Match(result.Outputs)
	})
	if err != nil {
		s.ErrC <- err
		return
	}
	for _, wf := range workflows {
		nextStep, err := wf.FindNode(wf.Trigger.NodeKey)
		if err != nil {
			s.ErrC <- err
			continue
		}
		if _, err := s.execution.Execute(wf.Hash, nextStep.InstanceHash, nil, result.Hash, wf.Trigger.NodeKey, nextStep.TaskKey, result.Outputs, []string{}); err != nil {
			s.ErrC <- err
		}
	}
}

func (s *Scheduler) workflowsMatchingFilter(filter func(wf *workflow.Workflow) bool) ([]*workflow.Workflow, error) {
	workflows, err := s.workflow.List()
	if err != nil {
		return nil, err
	}
	wfs := make([]*workflow.Workflow, 0)
	for _, wf := range workflows {
		if filter(wf) {
			wfs = append(wfs, wf)
		}
	}
	return wfs, nil
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
	for _, edge := range wf.EdgesFrom(exec.StepID) {
		inputs, err := s.mapInputs(wf.Hash, exec, edge)
		if err != nil {
			s.ErrC <- err
			return
		}
		nextStep, err := wf.FindNode(edge.Dst)
		if err != nil {
			s.ErrC <- err
			continue
		}
		if _, err := s.execution.Execute(wf.Hash, nextStep.InstanceHash, nil, exec.Hash, edge.Dst, nextStep.TaskKey, inputs, []string{}); err != nil {
			s.ErrC <- err
		}
	}
}

func (s *Scheduler) mapInputs(wfHash hash.Hash, prevExec *execution.Execution, edge workflow.Edge) (map[string]interface{}, error) {
	if len(edge.Inputs) == 0 {
		return prevExec.Outputs, nil
	}
	inputs := make(map[string]interface{})
	for _, input := range edge.Inputs {
		value, err := s.resolveInput(wfHash, prevExec, input.Ref.NodeKey, input.Ref.Key)
		if err != nil {
			return nil, err
		}
		inputs[input.Key] = value
	}
	return inputs, nil
}

func (s *Scheduler) resolveInput(wfHash hash.Hash, exec *execution.Execution, nodeKey string, outputKey string) (interface{}, error) {
	if !wfHash.Equal(exec.WorkflowHash) {
		return nil, fmt.Errorf("reference's nodeKey not found")
	}
	if exec.StepID != nodeKey {
		parent, err := s.execution.Get(exec.ParentHash)
		if err != nil {
			return nil, err
		}
		return s.resolveInput(wfHash, parent, nodeKey, outputKey)
	}
	return exec.Outputs[outputKey], nil
}
