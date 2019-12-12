package orchestrator

import (
	"context"
	"errors"
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
func New(event EventSDK, execution ExecutionSDK, process ProcessSDK, runner RunnerSDK) *Orchestrator {
	return &Orchestrator{
		event:     event,
		execution: execution,
		process:   process,
		runner:    runner,
		ErrC:      make(chan error),
	}
}

// Start the process engine
func (s *Orchestrator) Start() error {
	if s.eventStream != nil || s.executionStream != nil {
		return fmt.Errorf("process orchestrator already running")
	}

	s.eventStream = s.event.GetStream(nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	executionStream, errC, err := s.execution.Stream(ctx, &api.StreamExecutionRequest{
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
		WithField("node.key", n.Key).
		WithField("type", fmt.Sprintf("%T", n)).Debug("process process")
	if task := n.GetTask(); task != nil {
		// This returns directly because a task cannot process its children.
		// Children will be processed only when the execution is done and the dependencies are resolved
		return s.processTask(n.Key, task, wf, exec, event, data)
	} else if m := n.GetMap(); m != nil {
		var err error
		data, err = s.processMap(m.Outputs, wf, exec, data)
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
		if err := s.executeNode(wf, children, exec, event, data); err != nil {
			// does not return an error to continue to process other tasks if needed
			s.ErrC <- err
		}
	}
	return nil
}

func (s *Orchestrator) processMap(outputs map[string]*process.Process_Node_Map_Output, wf *process.Process, exec *execution.Execution, data *types.Struct) (*types.Struct, error) {
	result := &types.Struct{
		Fields: make(map[string]*types.Value),
	}
	for key, output := range outputs {
		value, err := s.outputToValue(output, wf, exec, data)
		if err != nil {
			return nil, err
		}
		result.Fields[key] = value
	}
	return result, nil
}

func (s *Orchestrator) outputToValue(output *process.Process_Node_Map_Output, wf *process.Process, exec *execution.Execution, data *types.Struct) (*types.Value, error) {
	switch v := output.GetValue().(type) {
	case *process.Process_Node_Map_Output_Null_:
		return &types.Value{Kind: &types.Value_NullValue{NullValue: types.NullValue_NULL_VALUE}}, nil
	case *process.Process_Node_Map_Output_StringConst:
		return &types.Value{Kind: &types.Value_StringValue{StringValue: v.StringConst}}, nil
	case *process.Process_Node_Map_Output_DoubleConst:
		return &types.Value{Kind: &types.Value_NumberValue{NumberValue: v.DoubleConst}}, nil
	case *process.Process_Node_Map_Output_BoolConst:
		return &types.Value{Kind: &types.Value_BoolValue{BoolValue: v.BoolConst}}, nil
	case *process.Process_Node_Map_Output_Map_:
		out, err := s.processMap(v.Map.Outputs, wf, exec, data)
		if err != nil {
			return nil, err
		}
		return &types.Value{Kind: &types.Value_StructValue{StructValue: out}}, nil
	case *process.Process_Node_Map_Output_List_:
		var values []*types.Value
		for i := range v.List.Outputs {
			value, err := s.outputToValue(v.List.Outputs[i], wf, exec, data)
			if err != nil {
				return nil, err
			}

			values = append(values, value)
		}
		return &types.Value{
			Kind: &types.Value_ListValue{
				ListValue: &types.ListValue{
					Values: values,
				},
			},
		}, nil
	case *process.Process_Node_Map_Output_Ref:
		node, err := wf.FindNode(v.Ref.NodeKey)
		if err != nil {
			return nil, err
		}
		if node.GetTask() != nil {
			return s.resolveInput(wf.Hash, exec, v.Ref.NodeKey, v.Ref.Key)
		}
		return data.Fields[v.Ref.Key], nil
	default:
		return nil, errors.New("unknown output")
	}
}

func (s *Orchestrator) resolveInput(wfHash hash.Hash, exec *execution.Execution, nodeKey string, outputKey string) (*types.Value, error) {
	if !wfHash.Equal(exec.ProcessHash) {
		return nil, fmt.Errorf("reference's nodeKey not found")
	}
	if exec.NodeKey != nodeKey {
		parent, err := s.execution.Get(exec.ParentHash)
		if err != nil {
			return nil, err
		}
		return s.resolveInput(wfHash, parent, nodeKey, outputKey)
	}
	return exec.Outputs.Fields[outputKey], nil
}

func (s *Orchestrator) processTask(nodeKey string, task *process.Process_Node_Task, wf *process.Process, exec *execution.Execution, event *event.Event, data *types.Struct) error {
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
		EventHash:    eventHash,
		ParentHash:   execHash,
		NodeKey:      nodeKey,
		TaskKey:      task.TaskKey,
		Inputs:       data,
		ExecutorHash: executor.Hash,
		Tags:         nil,
	})
	return err
}
