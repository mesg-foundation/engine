package orchestrator

import (
	"fmt"
	"math/rand"

	"github.com/mesg-foundation/engine/event"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/protobuf/types"
	runnersdk "github.com/mesg-foundation/engine/sdk/runner"
	"github.com/sirupsen/logrus"
)

// New creates a new Process instance
func New(event EventSDK, execution ExecutionSDK, process ProcessSDK, runner RunnerSDK, accountName, accountPassword string) *Orchestrator {
	return &Orchestrator{
		event:           event,
		execution:       execution,
		process:         process,
		runner:          runner,
		ErrC:            make(chan error),
		accountName:     accountName,
		accountPassword: accountPassword,
	}
}

// Start the process engine
func (s *Orchestrator) Start() error {
	if s.eventStream != nil || s.executionStream != nil {
		return fmt.Errorf("process orchestrator already running")
	}
	s.eventStream = s.event.GetStream(nil)
	executionStream, err := s.execution.Stream(&api.StreamExecutionRequest{
		Filter: &api.StreamExecutionRequest_Filter{
			Statuses: []execution.Status{execution.Status_Completed},
		},
	})
	if err != nil {
		return err
	}
	s.executionStream = executionStream
	for {
		select {
		case event := <-s.eventStream.C:
			go s.execute(s.eventFilter(event), nil, event, event.Data)
		case execution := <-s.executionStream:
			go s.execute(s.resultFilter(execution), execution, nil, execution.Outputs)
			go s.execute(s.dependencyFilter(execution), execution, nil, execution.Outputs)
		}
	}
	close(executionStream)
	return nil
}

func (s *Orchestrator) eventFilter(event *event.Event) func(wf *process.Process, node *process.Process_Node) (bool, error) {
	return func(wf *process.Process, node *process.Process_Node) (bool, error) {
		if e := node.GetEvent(); e != nil {
			return e.InstanceHash.Equal(event.InstanceHash) && e.EventKey == event.Key, nil
		}
		return false, nil
	}
}

func (s *Orchestrator) resultFilter(exec *execution.Execution) func(wf *process.Process, node *process.Process_Node) (bool, error) {
	return func(wf *process.Process, node *process.Process_Node) (bool, error) {
		if result := node.GetResult(); result != nil {
			return result.InstanceHash.Equal(exec.InstanceHash) && result.TaskKey == exec.TaskKey, nil
		}
		return false, nil
	}
}

func (s *Orchestrator) dependencyFilter(exec *execution.Execution) func(wf *process.Process, node *process.Process_Node) (bool, error) {
	return func(wf *process.Process, node *process.Process_Node) (bool, error) {
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

func (s *Orchestrator) findNodes(wf *process.Process, filter func(wf *process.Process, n *process.Process_Node) (bool, error)) []*process.Process_Node {
	return wf.FindNodes(func(n *process.Process_Node) bool {
		res, err := filter(wf, n)
		if err != nil {
			s.ErrC <- err
		}
		return res
	})
}

func (s *Orchestrator) execute(filter func(wf *process.Process, node *process.Process_Node) (bool, error), exec *execution.Execution, event *event.Event, data *types.Struct) {
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

func (s *Orchestrator) executeNode(wf *process.Process, n *process.Process_Node, exec *execution.Execution, event *event.Event, data *types.Struct) error {
	logrus.WithField("module", "orchestrator").
		WithField("nodeID", n.ID()).
		WithField("type", fmt.Sprintf("%T", n)).Debug("process process")
	if task := n.GetTask(); task != nil {
		// This returns directly because a task cannot process its children.
		// Children will be processed only when the execution is done and the dependencies are resolved
		return s.processTask(task, wf, exec, event, data)
	} else if m := n.GetMap(); m != nil {
		var err error
		data, err = s.processMap(m, wf, exec, data)
		if err != nil {
			return err
		}
	} else if filter := n.GetFilter(); filter != nil {
		if !filter.Match(data) {
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

func (s *Orchestrator) processMap(mapping *process.Process_Node_Map, wf *process.Process, exec *execution.Execution, data *types.Struct) (*types.Struct, error) {
	result := &types.Struct{
		Fields: make(map[string]*types.Value),
	}
	for _, output := range mapping.Outputs {
		if ref := output.GetRef(); ref != nil {
			node, err := wf.FindNode(ref.NodeKey)
			if err != nil {
				return nil, err
			}
			if node.GetTask() != nil {
				value, err := s.resolveInput(wf.Hash, exec, ref.NodeKey, ref.Key)
				if err != nil {
					return nil, err
				}
				result.Fields[output.Key] = value
			} else {
				result.Fields[output.Key] = data.Fields[ref.Key]
			}
		} else if constant := output.GetConstant(); constant != nil {
			result.Fields[output.Key] = constant
		}
	}
	return result, nil
}

func (s *Orchestrator) resolveInput(wfHash hash.Hash, exec *execution.Execution, nodeKey string, outputKey string) (*types.Value, error) {
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
	return exec.Outputs.Fields[outputKey], nil
}

func (s *Orchestrator) processTask(task *process.Process_Node_Task, wf *process.Process, exec *execution.Execution, event *event.Event, data *types.Struct) error {
	var eventHash, execHash hash.Hash
	if event != nil {
		eventHash = event.Hash
	}
	if exec != nil {
		execHash = exec.Hash
	}
	executors, err := s.runner.List(&runnersdk.Filter{
		InstanceHash: task.InstanceHash,
	})
	if err != nil {
		return err
	}
	if len(executors) == 0 {
		return fmt.Errorf("no runner is running instance %q", task.InstanceHash)
	}
	executor := executors[rand.Intn(len(executors))]
	_, err = s.execution.Create(&api.CreateExecutionRequest{
		ProcessHash:  wf.Hash,
		InstanceHash: task.InstanceHash,
		EventHash:    eventHash,
		ParentHash:   execHash,
		StepID:       task.Key,
		TaskKey:      task.TaskKey,
		Inputs:       data,
		ExecutorHash: executor.Hash,
		Tags:         nil,
	}, s.accountName, s.accountPassword)
	return err
}
