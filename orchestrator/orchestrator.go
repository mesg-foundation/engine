package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/mesg-foundation/engine/event"
	"github.com/mesg-foundation/engine/event/publisher"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/ext/xrand"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/runner"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

func init() {
	xrand.SeedInit()
}

// Store is an interface for fetching all data the orchestrator needs.
type Store interface {
	// FetchProcesses returns all processes.
	FetchProcesses(ctx context.Context) ([]*process.Process, error)

	// FetchExecution returns one execution from its hash.
	FetchExecution(ctx context.Context, hash hash.Hash) (*execution.Execution, error)

	// FetchRunners returns all runners of an instance.
	FetchRunners(ctx context.Context, instanceHash hash.Hash) ([]*runner.Runner, error)

	// CreateExecution creates an execution.
	CreateExecution(ctx context.Context, taskKey string, inputs *types.Struct, tags []string, parentHash hash.Hash, eventHash hash.Hash, processHash hash.Hash, nodeKey string, executorHash hash.Hash) (hash.Hash, error)

	// SubscribeToNewCompletedExecutions returns a chan that will contain newly completed execution.
	SubscribeToNewCompletedExecutions(ctx context.Context) (<-chan *execution.Execution, error)
}

// Orchestrator manages the executions based on the definition of the processes
type Orchestrator struct {
	store  Store
	ep     *publisher.EventPublisher
	logger tmlog.Logger

	eventStream     *event.Listener
	executionStream <-chan *execution.Execution
	stopC           chan bool
}

// New creates a new Process instance
func New(store Store, ep *publisher.EventPublisher, logger tmlog.Logger) *Orchestrator {
	return &Orchestrator{
		store:  store,
		ep:     ep,
		logger: logger.With("module", "orchestrator"),

		stopC: make(chan bool),
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
	if err := s.startExecutionStream(ctx); err != nil {
		return err
	}

	for {
		select {
		case event := <-s.eventStream.C:
			go s.execute(s.eventFilter(event), nil, event, event.Data)
		case execution := <-s.executionStream:
			go s.execute(s.resultFilter(execution), execution, nil, execution.Outputs)
			go s.execute(s.dependencyFilter(execution), execution, nil, execution.Outputs)
		case <-s.stopC:
			return nil
		}
	}
}

// Stop stops the orchestrator engine
func (s *Orchestrator) Stop() {
	s.stopC <- true
	s.eventStream.Close()
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
			s.logger.With(keyvals(wf, n, nil, nil, nil)...).Error(err.Error())
		}
		return res
	})
}

func (s *Orchestrator) execute(filter func(wf *process.Process, node *process.Process_Node) (bool, error), exec *execution.Execution, event *event.Event, data *types.Struct) {
	processes, err := s.store.FetchProcesses(context.Background())
	if err != nil {
		s.logger.Error(err.Error())
		return
	}
	for _, wf := range processes {
		for _, node := range s.findNodes(wf, filter) {
			if err := s.executeNode(wf, node, exec, event, data); err != nil {
				s.logger.With(keyvals(wf, node, exec, event, data)...).Error(err.Error())
			}
		}
	}
}

func (s *Orchestrator) executeNode(wf *process.Process, n *process.Process_Node, exec *execution.Execution, event *event.Event, data *types.Struct) error {
	log := s.logger.With(keyvals(wf, n, exec, event, data)...)
	// Process the node
	switch x := n.Type.(type) {
	case *process.Process_Node_Task_:
		// This returns directly because a task cannot process its children.
		// Children will be processed only when the execution is done and the dependencies are resolved
		createdExecHash, err := s.processTask(n, x.Task, wf, exec, event, data)
		if err != nil {
			return err
		}
		log.With("createdExecHash", createdExecHash.String()).Info("execution created")
		return nil // stop workflow execution
	case *process.Process_Node_Map_:
		var err error
		data, err = s.processMap(n.Key, x.Map.Outputs, wf, exec, data)
		if err != nil {
			return err
		}
		if result, err := json.Marshal(data); err == nil {
			log = log.With("output", string(result))
		}
	case *process.Process_Node_Filter_:
		if result, err := json.Marshal(x.Filter); err == nil {
			log = log.With("filter", string(result))
		}
		if !s.filterMatch(x.Filter, wf, n, exec, data) {
			log.Info("filter does not match data")
			return nil // stop workflow execution
		}
	}

	// Process the children
	for _, childrenID := range wf.ChildrenKeys(n.Key) {
		log = log.With("to", childrenID)
		children, err := wf.FindNode(childrenID)
		if err != nil {
			// does not return an error to continue to process other tasks if needed
			log.Error(err.Error())
			continue
		}
		log.Info("executed process transition")
		if err := s.executeNode(wf, children, exec, event, data); err != nil {
			// does not return an error to continue to process other tasks if needed
			s.logger.With(keyvals(wf, children, exec, event, data)...).Error(err.Error())

		}
	}

	return nil
}

// filterMatch returns true if the data match the current list of filters.
func (s *Orchestrator) filterMatch(f *process.Process_Node_Filter, wf *process.Process, n *process.Process_Node, exec *execution.Execution, data *types.Struct) bool {
	log := s.logger.With(keyvals(wf, n, exec, nil, data)...)
	for _, condition := range f.Conditions {
		resolvedData, err := s.resolveRef(wf, exec, n.Key, data, condition.Ref)
		if err != nil {
			log.Error(err.Error())
			return false
		}
		match, err := condition.Match(resolvedData)
		if err != nil {
			log.Error(err.Error())
		}
		if !match {
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
		return nil, fmt.Errorf("unknown output")
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
	parentExec, err := s.store.FetchExecution(context.Background(), exec.ParentHash)
	if err != nil {
		return nil, err
	}
	return s.resolveInput(wfHash, parentExec, refNodeKey, path)
}

// processTask create the request to execute the task.
func (s *Orchestrator) processTask(node *process.Process_Node, task *process.Process_Node_Task, wf *process.Process, exec *execution.Execution, event *event.Event, data *types.Struct) (hash.Hash, error) {
	var eventHash, execHash hash.Hash
	if event != nil {
		eventHash = event.Hash
	}
	if exec != nil {
		execHash = exec.Hash
	}
	executors, err := s.store.FetchRunners(context.Background(), task.InstanceHash)
	if err != nil {
		return nil, err
	}
	if len(executors) == 0 {
		return nil, fmt.Errorf("no runner is running instance %q", task.InstanceHash)
	}
	executor := executors[rand.Intn(len(executors))]

	// create execution
	return s.store.CreateExecution(
		context.Background(),
		task.TaskKey,
		data,
		nil,
		execHash,
		eventHash,
		wf.Hash,
		node.Key,
		executor.Hash,
	)
}

// startExecutionStream returns execution that matches given hash.
func (s *Orchestrator) startExecutionStream(ctx context.Context) error {
	var err error
	s.executionStream, err = s.store.SubscribeToNewCompletedExecutions(ctx)
	return err
}

func keyvals(proc *process.Process, node *process.Process_Node, parentExec *execution.Execution, event *event.Event, data *types.Struct) []interface{} {
	keyvals := []interface{}{}
	if proc != nil {
		keyvals = append(keyvals, "processHash", proc.Hash.String())
	}
	if node != nil {
		keyvals = append(keyvals,
			"from", node.Key,
			"type", node.TypeString(),
		)
	}
	if event != nil {
		keyvals = append(keyvals, "eventHash", event.Hash.String())
	}
	if parentExec != nil {
		keyvals = append(keyvals, "parentHash", parentExec.Hash.String())
	}
	if data != nil {
		if result, err := json.Marshal(data); err == nil {
			keyvals = append(keyvals, "input", string(result))
		}
	}
	return keyvals
}
