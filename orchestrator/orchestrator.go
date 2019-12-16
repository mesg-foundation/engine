package orchestrator

import (
	"context"
	"errors"
	"fmt"
	"log"
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
		return parents[0] == exec.RefKey, nil
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
		data, err = s.processMap(m.Outputs, wf, exec, event, data)
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

func (s *Orchestrator) processMap(outputs map[string]*process.Process_Node_Map_Output, wf *process.Process, exec *execution.Execution, event *event.Event, data *types.Struct) (*types.Struct, error) {
	result := &types.Struct{
		Fields: make(map[string]*types.Value),
	}
	for key, output := range outputs {
		value, err := s.outputToValue(output, wf, exec, event, data)
		if err != nil {
			return nil, err
		}
		result.Fields[key] = value
	}
	return result, nil
}

func (s *Orchestrator) outputToValue(output *process.Process_Node_Map_Output, wf *process.Process, exec *execution.Execution, event *event.Event, data *types.Struct) (*types.Value, error) {
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
		out, err := s.processMap(v.Map.Outputs, wf, exec, event, data)
		if err != nil {
			return nil, err
		}
		return &types.Value{Kind: &types.Value_StructValue{StructValue: out}}, nil
	case *process.Process_Node_Map_Output_List_:
		var values []*types.Value
		for i := range v.List.Outputs {
			value, err := s.outputToValue(v.List.Outputs[i], wf, exec, event, data)
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
		// XXX: HERE: we need previous node to compare key == refkey (we only have a event and exec)
		// if event != nil && event.Key == v.Ref.RefKey {
		// 	return resolveRef(event.Data, v.Ref.Path)
		// }
		// if exec != nil && exec.RefKey == v.Ref.RefKey {
		// 	return resolveRef(exec.Outputs, v.Ref.Path)
		// }

		if v.Ref.RefKey == "" {
			return resolveRef(data, v.Ref.Path)
		}
		// TODO: the following code is only for catching checking the node exist.
		nodes := wf.FindNodes(func(n *process.Process_Node) bool {
			task := n.GetTask()
			log.Println("getting task", n.Key)
			if task != nil {
				log.Println("found task", n.Key, task.RefKey, v.Ref.RefKey, task != nil && task.RefKey == v.Ref.RefKey)
			}
			return task != nil && task.RefKey == v.Ref.RefKey
		})
		log.Println("found n nodes", len(nodes))
		if len(nodes) != 1 {
			return nil, fmt.Errorf("node with reference's key %s not found", v.Ref.RefKey)
		}
		return s.resolveInput(wf.Hash, exec, v.Ref.RefKey, v.Ref.Path)
	default:
		return nil, errors.New("unknown output")
	}
}

func (s *Orchestrator) resolveInput(wfHash hash.Hash, exec *execution.Execution, refKey string, path *process.Process_Node_Map_Output_Reference_Path) (*types.Value, error) {
	if !wfHash.Equal(exec.ProcessHash) {
		return nil, fmt.Errorf("process with reference's key not found")
	}
	if exec.RefKey != refKey {
		parent, err := s.execution.Get(exec.ParentHash)
		if err != nil {
			return nil, err
		}
		return s.resolveInput(wfHash, parent, refKey, path)
	}
	return resolveRef(exec.Outputs, path)
}

func (s *Orchestrator) processTask(refKey string, task *process.Process_Node_Task, wf *process.Process, exec *execution.Execution, event *event.Event, data *types.Struct) error {
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
		RefKey:       refKey,
		TaskKey:      task.TaskKey,
		Inputs:       data,
		ExecutorHash: executor.Hash,
		Tags:         nil,
	})
	return err
}

func resolveRef(data *types.Struct, path *process.Process_Node_Map_Output_Reference_Path) (*types.Value, error) {
	if path == nil {
		return &types.Value{Kind: &types.Value_StructValue{StructValue: data}}, nil
	}

	var v *types.Value
	key, ok := path.Selector.(*process.Process_Node_Map_Output_Reference_Path_Key)
	if !ok {
		return nil, fmt.Errorf("orchestrator: first selector in the path must be a key")
	}

	v, ok = data.Fields[key.Key]
	if !ok {
		return nil, fmt.Errorf("orchestrator: key %s not found", key.Key)
	}

	for p := path.Path; p != nil; p = p.Path {
		switch s := p.Selector.(type) {
		case *process.Process_Node_Map_Output_Reference_Path_Key:
			str, ok := v.GetKind().(*types.Value_StructValue)
			if !ok {
				return nil, fmt.Errorf("orchestrator: can't get key from non-struct value")
			}
			if str.StructValue.GetFields() == nil {
				return nil, fmt.Errorf("orchestrator: can't get key from nil-struct")
			}
			v, ok = str.StructValue.Fields[s.Key]
			if !ok {
				return nil, fmt.Errorf("orchestrator: key %s not found", s.Key)
			}
		case *process.Process_Node_Map_Output_Reference_Path_Index:
			list, ok := v.GetKind().(*types.Value_ListValue)
			if !ok {
				return nil, fmt.Errorf("orchestrator: can't get index from non-list value")
			}

			if len(list.ListValue.GetValues()) <= int(s.Index) {
				return nil, fmt.Errorf("orchestrator: index %d out of range", s.Index)
			}
			v = list.ListValue.Values[s.Index]
		default:
			return nil, fmt.Errorf("orchestrator: unknown selector type %T", v)
		}
	}

	return v, nil
}
