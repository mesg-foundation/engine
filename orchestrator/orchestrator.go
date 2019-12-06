package orchestrator

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/mesg-foundation/engine/event"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/result"
	runnersdk "github.com/mesg-foundation/engine/sdk/runner"
	"github.com/sirupsen/logrus"
)

// New creates a new Process instance
func New(event EventSDK, exec ExecutionSDK, res ResultSDK, process ProcessSDK, runner RunnerSDK, accountName, accountPassword string) *Orchestrator {
	return &Orchestrator{
		event:           event,
		execution:       exec,
		result:          res,
		process:         process,
		runner:          runner,
		ErrC:            make(chan error),
		accountName:     accountName,
		accountPassword: accountPassword,
	}
}

// Start the process engine
func (s *Orchestrator) Start() error {
	if s.eventStream != nil || s.resultStream != nil {
		return fmt.Errorf("process orchestrator already running")
	}

	s.eventStream = s.event.GetStream(nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	resultStream, errC, err := s.result.Stream(ctx, &api.StreamResultRequest{})
	if err != nil {
		return err
	}
	s.resultStream = resultStream
	for {
		select {
		case event := <-s.eventStream.C:
			go s.execute(s.eventFilter(event), nil, nil, event, event.Data)
		case res := <-s.resultStream:
			go func(res *result.Result) {
				if x, ok := res.GetResult().(*result.Result_Error); ok && x != nil {
					// discard result containing error
					return
				}
				exec, err := s.execution.Get(res.ExecutionHash)
				if err != nil {
					s.ErrC <- err
					return
				}
				go s.execute(s.resultFilter(exec), exec, res, nil, res.GetOutputs())
				go s.execute(s.dependencyFilter(exec), exec, res, nil, res.GetOutputs())
			}(res)
		case err := <-errC:
			s.ErrC <- err
		}
	}
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
		parents := wf.ParentKeys(node.Key)
		if len(parents) == 0 {
			return false, nil
		}
		if len(parents) > 1 {
			return false, fmt.Errorf("multi parents not supported")
		}
		return parents[0] == exec.NodeKey, nil
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

func (s *Orchestrator) execute(filter func(wf *process.Process, node *process.Process_Node) (bool, error), exec *execution.Execution, res *result.Result, event *event.Event, data *types.Struct) {
	processes, err := s.process.List()
	if err != nil {
		s.ErrC <- err
		return
	}
	for _, wf := range processes {
		for _, node := range s.findNodes(wf, filter) {
			if err := s.executeNode(wf, node, exec, res, event, data); err != nil {
				s.ErrC <- err
			}
		}
	}
}

func (s *Orchestrator) executeNode(wf *process.Process, n *process.Process_Node, exec *execution.Execution, res *result.Result, event *event.Event, data *types.Struct) error {
	logrus.WithField("module", "orchestrator").
		WithField("node.key", n.Key).
		WithField("type", fmt.Sprintf("%T", n)).Debug("process process")
	if task := n.GetTask(); task != nil {
		// This returns directly because a task cannot process its children.
		// Children will be processed only when the execution is done and the dependencies are resolved
		return s.processTask(n.Key, task, wf, res, event, data)
	} else if m := n.GetMap(); m != nil {
		var err error
		data, err = s.processMap(m, wf, exec, res, data)
		if err != nil {
			return err
		}
	} else if filter := n.GetFilter(); filter != nil {
		if !filter.Match(data) {
			return nil
		}
	}
	for _, childrenID := range wf.ChildrenKeys(n.Key) {
		children, err := wf.FindNode(childrenID)
		if err != nil {
			// does not return an error to continue to process other tasks if needed
			s.ErrC <- err
			continue
		}
		if err := s.executeNode(wf, children, exec, res, event, data); err != nil {
			// does not return an error to continue to process other tasks if needed
			s.ErrC <- err
		}
	}
	return nil
}

func (s *Orchestrator) processMap(mapping *process.Process_Node_Map, wf *process.Process, exec *execution.Execution, res *result.Result, data *types.Struct) (*types.Struct, error) {
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
				value, err := s.resolveInput(wf.Hash, exec, res, ref.NodeKey, ref.Key)
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

func (s *Orchestrator) resolveInput(wfHash hash.Hash, exec *execution.Execution, res *result.Result, nodeKey string, outputKey string) (*types.Value, error) {
	if !wfHash.Equal(exec.ProcessHash) {
		return nil, fmt.Errorf("reference's nodeKey not found")
	}
	if exec.NodeKey != nodeKey {
		resParent, err := s.result.Get(exec.ParentResultHash)
		if err != nil {
			return nil, err
		}
		execParent, err := s.execution.Get(res.Hash)
		if err != nil {
			return nil, err
		}
		return s.resolveInput(wfHash, execParent, resParent, nodeKey, outputKey)
	}
	return res.GetOutputs().Fields[outputKey], nil
}

func (s *Orchestrator) processTask(nodeKey string, task *process.Process_Node_Task, wf *process.Process, res *result.Result, event *event.Event, data *types.Struct) error {
	var eventHash, resHash hash.Hash
	if event != nil {
		eventHash = event.Hash
	}
	if res != nil {
		resHash = res.Hash
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
		ProcessHash:      wf.Hash,
		EventHash:        eventHash,
		ParentResultHash: resHash,
		NodeKey:          nodeKey,
		TaskKey:          task.TaskKey,
		Inputs:           data,
		ExecutorHash:     executor.Hash,
		Tags:             nil,
	}, s.accountName, s.accountPassword)
	return err
}
