package orchestrator

import (
	"fmt"

	"github.com/mesg-foundation/engine/event"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/process"
	eventsdk "github.com/mesg-foundation/engine/sdk/event"
	executionsdk "github.com/mesg-foundation/engine/sdk/execution"
	processsdk "github.com/mesg-foundation/engine/sdk/process"
	"github.com/sirupsen/logrus"
)

// Orchestrator manages the executions based on the definition of the processes
type Orchestrator struct {
	event       *eventsdk.Event
	eventStream *eventsdk.Listener

	execution       *executionsdk.Execution
	executionStream *executionsdk.Listener

	process *processsdk.Process

	ErrC chan error
}

// New creates a new Process instance
func New(event *eventsdk.Event, execution *executionsdk.Execution, process *processsdk.Process) *Orchestrator {
	return &Orchestrator{
		event:     event,
		execution: execution,
		process:   process,
		ErrC:      make(chan error),
	}
}

// Start the process engine
func (o *Orchestrator) Start() error {
	if o.eventStream != nil || o.executionStream != nil {
		return fmt.Errorf("process orchestrator already running")
	}
	o.eventStream = o.event.GetStream(nil)
	o.executionStream = o.execution.GetStream(&executionsdk.Filter{
		Statuses: []execution.Status{execution.Completed},
	})
	for {
		select {
		case event := <-o.eventStream.C:
			go o.execute(o.eventFilter(event), nil, event, event.Data)
		case execution := <-o.executionStream.C:
			go o.execute(o.resultFilter(execution), execution, nil, execution.Outputs)
			go o.execute(o.dependencyFilter(execution), execution, nil, execution.Outputs)
		}
	}
}

func (o *Orchestrator) eventFilter(event *event.Event) func(wf *process.Process, node process.Node) (bool, error) {
	return func(wf *process.Process, node process.Node) (bool, error) {
		switch n := node.(type) {
		case *process.Event:
			return n.InstanceHash.Equal(event.InstanceHash) && n.EventKey == event.Key, nil
		default:
			return false, nil
		}
	}
}

func (o *Orchestrator) resultFilter(exec *execution.Execution) func(wf *process.Process, node process.Node) (bool, error) {
	return func(wf *process.Process, node process.Node) (bool, error) {
		switch n := node.(type) {
		case *process.Result:
			return n.InstanceHash.Equal(exec.InstanceHash) && n.TaskKey == exec.TaskKey, nil
		default:
			return false, nil
		}
	}
}

func (o *Orchestrator) dependencyFilter(exec *execution.Execution) func(wf *process.Process, node process.Node) (bool, error) {
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

func (o *Orchestrator) findNodes(wf *process.Process, filter func(wf *process.Process, n process.Node) (bool, error)) []process.Node {
	return wf.FindNodes(func(n process.Node) bool {
		res, err := filter(wf, n)
		if err != nil {
			o.ErrC <- err
		}
		return res
	})
}

func (o *Orchestrator) execute(filter func(wf *process.Process, node process.Node) (bool, error), exec *execution.Execution, event *event.Event, data map[string]interface{}) {
	processes, err := o.process.List()
	if err != nil {
		o.ErrC <- err
		return
	}
	for _, wf := range processes {
		for _, node := range o.findNodes(wf, filter) {
			if err := o.executeNode(wf, node, exec, event, data); err != nil {
				o.ErrC <- err
			}
		}
	}
}

func (o *Orchestrator) executeNode(wf *process.Process, n process.Node, exec *execution.Execution, event *event.Event, data map[string]interface{}) error {
	logrus.WithField("module", "orchestrator").WithField("nodeID", n.ID()).WithField("type", fmt.Sprintf("%T", n)).Debug("process process")
	if node, ok := n.(*process.Task); ok {
		// This returns directly because a task cannot process its children.
		// Children will be processed only when the execution is done and the dependencies are resolved
		return o.processTask(node, wf, exec, event, data)
	}
	if node, ok := n.(*process.Map); ok {
		var err error
		data, err = o.processMap(node, wf, exec, data)
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
			o.ErrC <- err
			continue
		}
		if err := o.executeNode(wf, children, exec, event, data); err != nil {
			// does not return an error to continue to process other tasks if needed
			o.ErrC <- err
		}
	}
	return nil
}

func (o *Orchestrator) processMap(mapping *process.Map, wf *process.Process, exec *execution.Execution, data map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for _, output := range mapping.Outputs {
		node, err := wf.FindNode(output.Ref.NodeKey)
		if err != nil {
			return nil, err
		}
		_, isTask := node.(*process.Task)
		if isTask {
			value, err := o.resolveInput(wf.Hash, exec, output.Ref.NodeKey, output.Ref.Key)
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

func (o *Orchestrator) resolveInput(wfHash hash.Hash, exec *execution.Execution, nodeKey string, outputKey string) (interface{}, error) {
	if !wfHash.Equal(exec.ProcessHash) {
		return nil, fmt.Errorf("reference's nodeKey not found")
	}
	if exec.StepID != nodeKey {
		parent, err := o.execution.Get(exec.ParentHash)
		if err != nil {
			return nil, err
		}
		return o.resolveInput(wfHash, parent, nodeKey, outputKey)
	}
	return exec.Outputs[outputKey], nil
}

func (o *Orchestrator) processTask(task *process.Task, wf *process.Process, exec *execution.Execution, event *event.Event, data map[string]interface{}) error {
	var eventHash, execHash hash.Hash
	if event != nil {
		eventHash = event.Hash
	}
	if exec != nil {
		execHash = exec.Hash
	}
	_, err := o.execution.Execute(wf.Hash, task.InstanceHash, eventHash, execHash, task.ID(), task.TaskKey, data, nil)
	return err
}
