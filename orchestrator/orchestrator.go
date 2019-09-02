package orchestrator

import (
	"fmt"

	"github.com/gogo/protobuf/types"
	"github.com/mesg-foundation/engine/event"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/protobuf/convert"
	eventsdk "github.com/mesg-foundation/engine/sdk/event"
	executionsdk "github.com/mesg-foundation/engine/sdk/execution"
	processesdk "github.com/mesg-foundation/engine/sdk/process"
	"github.com/sirupsen/logrus"
)

// Orchestrator manages the executions based on the definition of the processes
type Orchestrator struct {
	event       *eventsdk.Event
	eventStream *eventsdk.Listener

	execution       *executionsdk.Execution
	executionStream *executionsdk.Listener

	process *processesdk.Process

	ErrC chan error
}

// New creates a new Process instance
func New(event *eventsdk.Event, execution *executionsdk.Execution, process *processesdk.Process) *Orchestrator {
	return &Orchestrator{
		event:     event,
		execution: execution,
		process:   process,
		ErrC:      make(chan error),
	}
}

// Start the process engine
func (s *Orchestrator) Start() error {
	if s.eventStream != nil || s.executionStream != nil {
		return fmt.Errorf("process orchestrator already running")
	}
	s.eventStream = s.event.GetStream(nil)
	s.executionStream = s.execution.GetStream(&executionsdk.Filter{
		Statuses: []execution.Status{execution.Status_Completed},
	})
	for {
		select {
		case event := <-s.eventStream.C:
			data := make(map[string]interface{})
			if err := convert.Marshal(event.Data, &data); err != nil {
				s.ErrC <- err
				continue
			}

			go s.execute(s.eventFilter(event), nil, event, data)
		case execution := <-s.executionStream.C:
			outputs := make(map[string]interface{})
			if err := convert.Marshal(execution.Outputs, &outputs); err != nil {
				s.ErrC <- err
				continue
			}
			go s.execute(s.resultFilter(execution), execution, nil, outputs)
			go s.execute(s.dependencyFilter(execution), execution, nil, outputs)
		}
	}
}

func (s *Orchestrator) eventFilter(event *event.Event) func(wf *process.Process, node process.Node) (bool, error) {
	return func(wf *process.Process, node process.Node) (bool, error) {
		switch n := node.(type) {
		case *process.Event:
			return n.InstanceHash.Equal(event.InstanceHash) && n.EventKey == event.Key, nil
		default:
			return false, nil
		}
	}
}

func (s *Orchestrator) resultFilter(exec *execution.Execution) func(wf *process.Process, node process.Node) (bool, error) {
	return func(wf *process.Process, node process.Node) (bool, error) {
		switch n := node.(type) {
		case *process.Result:
			return n.InstanceHash.Equal(exec.InstanceHash) && n.TaskKey == exec.TaskKey, nil
		default:
			return false, nil
		}
	}
}

func (s *Orchestrator) dependencyFilter(exec *execution.Execution) func(wf *process.Process, node process.Node) (bool, error) {
	return func(wf *process.Process, node process.Node) (bool, error) {
		if !exec.ProcessHash.Equal(wf.Hash) {
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

func (s *Orchestrator) findNodes(wf *process.Process, filter func(wf *process.Process, n process.Node) (bool, error)) []process.Node {
	return wf.FindNodes(func(n process.Node) bool {
		res, err := filter(wf, n)
		if err != nil {
			s.ErrC <- err
		}
		return res
	})
}

func (s *Orchestrator) execute(filter func(wf *process.Process, node process.Node) (bool, error), exec *execution.Execution, event *event.Event, data map[string]interface{}) {
	processes, err := s.process.List()
	if err != nil {
		s.ErrC <- err
		return
	}
	for _, wf := range processes {
		for _, node := range s.findNodes(wf, filter) {
			if err := s.executeNode(wf, node, exec, event, data); err != nil {
				s.ErrC <- err
			}
		}
	}
}

func (s *Orchestrator) executeNode(wf *process.Process, n process.Node, exec *execution.Execution, event *event.Event, data map[string]interface{}) error {
	logrus.WithField("module", "orchestrator").WithField("nodeID", n.ID()).WithField("type", fmt.Sprintf("%T", n)).Debug("process process")
	if node, ok := n.(*process.Task); ok {
		// This returns directly because a task cannot process its children.
		// Children will be processed only when the execution is done and the dependencies are resolved
		return s.processTask(node, wf, exec, event, data)
	}
	if node, ok := n.(*process.Map); ok {
		var err error
		data, err = s.processMap(node, wf, exec, data)
		if err != nil {
			return err
		}
	}
	if node, ok := n.(*process.Filter); ok {
		if !node.Filter.Match(data) {
			return nil
		}
	}
	for _, childrenID := range wf.ChildrenIDs(n.ID()) {
		children, err := wf.FindNode(childrenID)
		if err != nil {
			// does not return an error to continue to process other tasks if needed
			s.ErrC <- err
			continue
		}
		if err := s.executeNode(wf, children, exec, event, data); err != nil {
			// does not return an error to continue to process other tasks if needed
			s.ErrC <- err
		}
	}
	return nil
}

func (s *Orchestrator) processMap(mapping *process.Map, wf *process.Process, exec *execution.Execution, data map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for _, output := range mapping.Outputs {
		node, err := wf.FindNode(output.Ref.NodeKey)
		if err != nil {
			return nil, err
		}
		_, isTask := node.(*process.Task)
		if isTask {
			value, err := s.resolveInput(wf.Hash, exec, output.Ref.NodeKey, output.Ref.Key)
			if err != nil {
				return nil, err
			}
			result[output.Key] = value
		} else {
			result[output.Key] = data[output.Ref.Key]
		}
	}
	return result, nil
}

func (s *Orchestrator) resolveInput(wfHash hash.Hash, exec *execution.Execution, nodeKey string, outputKey string) (interface{}, error) {
	if !wfHash.Equal(exec.ProcessHash) {
		return nil, fmt.Errorf("reference's nodeKey not found")
	}
	if exec.StepID != nodeKey {
		parent, err := s.execution.Get(exec.ParentHash)
		if err != nil {
			return nil, err
		}
		return s.resolveInput(wfHash, parent, nodeKey, outputKey)
	}

	outputs := make(map[string]interface{})
	if err := convert.Marshal(exec.Outputs, &outputs); err != nil {
		return nil, err
	}

	return outputs[outputKey], nil
}

func (s *Orchestrator) processTask(task *process.Task, wf *process.Process, exec *execution.Execution, event *event.Event, data map[string]interface{}) error {
	var eventHash, execHash hash.Hash
	if event != nil {
		eventHash = event.Hash
	}
	if exec != nil {
		execHash = exec.Hash
	}

	execData := &types.Struct{}
	if err := convert.Unmarshal(data, execData); err != nil {
		return err
	}

	_, err := s.execution.Execute(wf.Hash, task.InstanceHash, eventHash, execHash, task.ID(), task.TaskKey, execData, nil)
	return err
}
