package orchestrator

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/event"
	"github.com/mesg-foundation/engine/event/publisher"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/sirupsen/logrus"
)

// New creates a new Process instance
func New(mc *cosmos.ModuleClient, ep *publisher.EventPublisher, execPrice string) *Orchestrator {
	return &Orchestrator{
		mc:        mc,
		ep:        ep,
		ErrC:      make(chan error),
		stopC:     make(chan bool),
		execPrice: execPrice,
	}
}

// Start the process engine
func (s *Orchestrator) Start() error {
	if s.eventStream != nil || s.executionStream != nil {
		return fmt.Errorf("process orchestrator already running")
	}

	s.eventStream = s.ep.GetStream(nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	executionStream, errC, err := s.mc.StreamExecution(ctx, &api.StreamExecutionRequest{
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
		case <-s.stopC:
			return nil
		}
	}
}

// Stop stops the orchestrator engine
func (s *Orchestrator) Stop() {
	s.stopC <- true
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
	processes, err := s.mc.ListProcess()
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
		data, err = s.processMap(n.Key, m.Outputs, wf, exec, data)
		if err != nil {
			return err
		}
	} else if filter := n.GetFilter(); filter != nil {
		if !s.filterMatch(filter, wf, n.Key, exec, data) {
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

// filterMatch returns true if the data match the current list of filters.
func (s *Orchestrator) filterMatch(f *process.Process_Node_Filter, wf *process.Process, nodeKey string, exec *execution.Execution, data *types.Struct) bool {
	for _, condition := range f.Conditions {
		resolvedData, err := s.resolveRef(wf, exec, nodeKey, data, condition.Ref)
		if err != nil {
			return false
		}
		if !condition.Match(resolvedData) {
			return false
		}
	}
	return true
}

// processMap constructs the Struct that the map indicates how to construct.
func (s *Orchestrator) processMap(nodeKey string, outputs map[string]*process.Process_Node_Map_Output, wf *process.Process, exec *execution.Execution, data *types.Struct) (*types.Struct, error) {
	result := &types.Struct{
		Fields: make(map[string]*types.Value),
	}
	for key, output := range outputs {
		value, err := s.outputToValue(nodeKey, output, wf, exec, data)
		if err != nil {
			return nil, err
		}
		result.Fields[key] = value
	}
	return result, nil
}

// outputToValue returns a specific value from an output.
func (s *Orchestrator) outputToValue(nodeKey string, output *process.Process_Node_Map_Output, wf *process.Process, exec *execution.Execution, data *types.Struct) (*types.Value, error) {
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
		out, err := s.processMap(nodeKey, v.Map.Outputs, wf, exec, data)
		if err != nil {
			return nil, err
		}
		return &types.Value{Kind: &types.Value_StructValue{StructValue: out}}, nil
	case *process.Process_Node_Map_Output_List_:
		var values []*types.Value
		for i := range v.List.Outputs {
			value, err := s.outputToValue(nodeKey, v.List.Outputs[i], wf, exec, data)
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
		return s.resolveRef(wf, exec, nodeKey, data, v.Ref)
	default:
		return nil, errors.New("unknown output")
	}
}

// resolveRef returns a specific value from a reference.
func (s *Orchestrator) resolveRef(wf *process.Process, exec *execution.Execution, nodeKey string, data *types.Struct, ref *process.Process_Node_Reference) (*types.Value, error) {
	refNode, err := wf.FindNode(ref.NodeKey)
	if err != nil {
		return nil, err
	}
	// if referenced node is a task, get its output
	if refNode.GetTask() != nil {
		return s.resolveInput(wf.Hash, exec, ref.NodeKey, ref.Path)
	}
	// check that the parent nodeKey == ref.NodeKey
	// this ensures that we can use directly the data of the previous node
	refToParent := false
	for _, parent := range wf.ParentKeys(nodeKey) {
		if parent == ref.NodeKey {
			refToParent = true
			break
		}
	}
	if !refToParent {
		return nil, fmt.Errorf("ref can only reference a parent node for non task nodes")
	}
	return ref.Path.Resolve(data)
}

// resolveInput returns a specific value from a reference path.
func (s *Orchestrator) resolveInput(wfHash hash.Hash, exec *execution.Execution, refNodeKey string, path *process.Process_Node_Reference_Path) (*types.Value, error) {
	// the wfHash condition only works for Task node, not for Result
	if !wfHash.Equal(exec.ProcessHash) {
		return nil, fmt.Errorf("reference's nodeKey not found")
	}
	// we reach the right execution, return it
	// but only works for Task as the execution related to the Result is not created by this process
	if exec.NodeKey == refNodeKey {
		return path.Resolve(exec.Outputs)
	}
	// get parentExec and do a recursive call
	parentExec, err := s.mc.GetExecution(exec.ParentHash)
	if err != nil {
		return nil, err
	}
	return s.resolveInput(wfHash, parentExec, refNodeKey, path)
}

// processTask create the request to execute the task.
func (s *Orchestrator) processTask(nodeKey string, task *process.Process_Node_Task, wf *process.Process, exec *execution.Execution, event *event.Event, data *types.Struct) error {
	var eventHash, execHash hash.Hash
	if event != nil {
		eventHash = event.Hash
	}
	if exec != nil {
		execHash = exec.Hash
	}
	executors, err := s.mc.ListRunner(&cosmos.FilterRunner{
		InstanceHash: task.InstanceHash,
	})
	if err != nil {
		return err
	}
	if len(executors) == 0 {
		return fmt.Errorf("no runner is running instance %q", task.InstanceHash)
	}
	executor := executors[rand.Intn(len(executors))]
	_, err = s.mc.CreateExecution(&api.CreateExecutionRequest{
		ProcessHash:  wf.Hash,
		EventHash:    eventHash,
		ParentHash:   execHash,
		NodeKey:      nodeKey,
		TaskKey:      task.TaskKey,
		Inputs:       data,
		ExecutorHash: executor.Hash,
		Price:        s.execPrice,
		Tags:         nil,
	})
	return err
}
