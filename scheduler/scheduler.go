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
	"github.com/sirupsen/logrus"
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
			go s.process(s.eventFilter(event), nil, event, event.Data)
		case execution := <-s.executionStream.C:
			go s.process(s.resultFilter(execution), execution, nil, execution.Outputs)
			go s.process(s.dependencyFilter(execution), execution, nil, execution.Outputs)
		}
	}
}

func (s *Scheduler) eventFilter(event *event.Event) func(wf *workflow.Workflow, node workflow.Node) (bool, error) {
	return func(wf *workflow.Workflow, node workflow.Node) (bool, error) {
		switch n := node.(type) {
		case *workflow.Event:
			return n.InstanceHash.Equal(event.InstanceHash) && n.EventKey == event.Key, nil
		default:
			return false, nil
	}
		}
		}

func (s *Scheduler) resultFilter(exec *execution.Execution) func(wf *workflow.Workflow, node workflow.Node) (bool, error) {
	return func(wf *workflow.Workflow, node workflow.Node) (bool, error) {
		switch n := node.(type) {
		case *workflow.Result:
			return n.InstanceHash.Equal(exec.InstanceHash) && n.TaskKey == exec.TaskKey, nil
		default:
			return false, nil
	}
}
}

func (s *Scheduler) dependencyFilter(exec *execution.Execution) func(wf *workflow.Workflow, node workflow.Node) (bool, error) {
	return func(wf *workflow.Workflow, node workflow.Node) (bool, error) {
		if !exec.WorkflowHash.Equal(wf.Hash) {
			return false, nil
		}
		parents := wf.ParentIDs(node.ID())
		if len(parents) == 0 {
			return false, nil
		}
		if len(parents) > 1 {
			return false, fmt.Errorf("multi parents not supported")
		}
		return parents[0] == exec.StepID, nil
	}
}

func (s *Scheduler) process(filter func(wf *workflow.Workflow, node workflow.Node) (bool, error), exec *execution.Execution, event *event.Event, data map[string]interface{}) {
	workflows, err := s.workflow.List()
	if err != nil {
		s.ErrC <- err
		return
	}
	for _, wf := range workflows {
		nodes := wf.FindNodes(func(n workflow.Node) bool {
			res, err := filter(wf, n)
		if err != nil {
			s.ErrC <- err
		}
			return res
		})
		for _, node := range nodes {
			if err := s.processNode(wf, node, exec, event, data); err != nil {
			s.ErrC <- err
		}
	}
}
}

func (s *Scheduler) processNode(wf *workflow.Workflow, n workflow.Node, exec *execution.Execution, event *event.Event, data map[string]interface{}) error {
	logrus.WithField("module", "orchestrator").WithField("nodeID", n.ID()).WithField("type", fmt.Sprintf("%T", n)).Debug("process workflow")
	switch node := n.(type) {
	case *workflow.Mapping:
		var err error
		data, err = s.processMapping(node, wf, exec, data)
	if err != nil {
			return err
	}
	case *workflow.Task:
		if err := s.processTask(node, wf, exec, event, data); err != nil {
			return err
		}
	}
	if s.canProcessChildren(n) {
		for _, childrenID := range wf.ChildrenIDs(n.ID()) {
			children, err := wf.FindNode(childrenID)
	if err != nil {
				// does not return an error to continue to process other tasks if needed
		s.ErrC <- err
				continue
	}
			if err := s.processNode(wf, children, exec, event, data); err != nil {
				// does not return an error to continue to process other tasks if needed
			s.ErrC <- err
			}
		}
		}
	return nil
		}

func (s *Scheduler) canProcessChildren(node workflow.Node) bool {
	switch node.(type) {
	case *workflow.Event:
		return true
	case *workflow.Result:
		return true
	case *workflow.Mapping:
		return true
	case *workflow.Task:
		return false
	}
	return false
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

func (s *Scheduler) processTask(task *workflow.Task, wf *workflow.Workflow, exec *execution.Execution, event *event.Event, data map[string]interface{}) error {
	var eventHash, execHash hash.Hash
	if event != nil {
		eventHash = event.Hash
	}
	if exec != nil {
		execHash = exec.Hash
	}
	_, err := s.execution.Execute(wf.Hash, task.InstanceHash, eventHash, execHash, task.ID(), task.TaskKey, data, []string{})
	return err
}
