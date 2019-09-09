package orchestrator

import (
	"testing"

	"github.com/mesg-foundation/engine/event"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/orchestrator/mocks"
	"github.com/mesg-foundation/engine/process"
	eventsdk "github.com/mesg-foundation/engine/sdk/event"
	executionsdk "github.com/mesg-foundation/engine/sdk/execution"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFilter(t *testing.T) {
	o := New(&mocks.EventSDK{}, &mocks.ExecutionSDK{}, &mocks.ProcessSDK{})
	p := process.Process{
		Hash: hash.Int(1),
		Nodes: []*process.Process_Node{
			&process.Process_Node{Type: &process.Process_Node_Event_{Event: &process.Process_Node_Event{
				Key:          "1",
				InstanceHash: hash.Int(1),
				EventKey:     "1",
			}}},
			&process.Process_Node{Type: &process.Process_Node_Task_{Task: &process.Process_Node_Task{
				InstanceHash: hash.Int(2),
				Key:          "2",
				TaskKey:      "2",
			}}},
			&process.Process_Node{Type: &process.Process_Node_Task_{Task: &process.Process_Node_Task{
				InstanceHash: hash.Int(3),
				Key:          "3",
				TaskKey:      "3",
			}}},
		},
		Edges: []*process.Process_Edge{
			{Src: "1", Dst: "2"},
			{Src: "2", Dst: "3"},
		},
	}
	var tests = []struct {
		filter func(wf *process.Process, node *process.Process_Node) (bool, error)
		p      *process.Process
		n      *process.Process_Node
		res    bool
		err    error
	}{
		{
			filter: o.eventFilter(&event.Event{InstanceHash: hash.Int(1), Key: "1"}),
			n:      p.Nodes[0],
			res:    true,
			err:    nil,
		},
		{
			filter: o.eventFilter(&event.Event{InstanceHash: hash.Int(1), Key: "2"}),
			n:      p.Nodes[0],
			res:    false,
			err:    nil,
		},
		{
			filter: o.eventFilter(&event.Event{InstanceHash: hash.Int(2), Key: "1"}),
			n:      p.Nodes[0],
			res:    false,
			err:    nil,
		},
		{
			filter: o.eventFilter(&event.Event{InstanceHash: hash.Int(2), Key: "1"}),
			n:      p.Nodes[1],
			res:    false,
			err:    nil,
		},
		{
			filter: o.resultFilter(&execution.Execution{InstanceHash: hash.Int(1), TaskKey: "1"}),
			n: &process.Process_Node{Type: &process.Process_Node_Result_{Result: &process.Process_Node_Result{
				InstanceHash: hash.Int(1),
				TaskKey:      "1",
			}}},
			res: true,
			err: nil,
		},
		{
			filter: o.resultFilter(&execution.Execution{InstanceHash: hash.Int(1), TaskKey: "1"}),
			n: &process.Process_Node{Type: &process.Process_Node_Result_{Result: &process.Process_Node_Result{
				InstanceHash: hash.Int(1),
				TaskKey:      "2",
			}}},
			res: false,
			err: nil,
		},
		{
			filter: o.resultFilter(&execution.Execution{InstanceHash: hash.Int(1), TaskKey: "1"}),
			n: &process.Process_Node{Type: &process.Process_Node_Result_{Result: &process.Process_Node_Result{
				InstanceHash: hash.Int(2),
				TaskKey:      "1",
			}}},
			res: false,
			err: nil,
		},
		{
			filter: o.resultFilter(&execution.Execution{InstanceHash: hash.Int(1), TaskKey: "1"}),
			n:      p.Nodes[0],
			res:    false,
			err:    nil,
		},
		{
			filter: o.dependencyFilter(&execution.Execution{InstanceHash: hash.Int(3), TaskKey: "2", ProcessHash: hash.Int(1), StepID: "2"}),
			p:      &p,
			n:      p.Nodes[2],
			res:    true,
			err:    nil,
		},
		{
			filter: o.dependencyFilter(&execution.Execution{InstanceHash: hash.Int(3), TaskKey: "2", ProcessHash: hash.Int(2), StepID: "2"}),
			p:      &p,
			n:      p.Nodes[2],
			res:    false,
			err:    nil,
		},
		{
			filter: o.dependencyFilter(&execution.Execution{InstanceHash: hash.Int(3), TaskKey: "2", ProcessHash: hash.Int(1), StepID: "1"}),
			p:      &p,
			n:      p.Nodes[2],
			res:    false,
			err:    nil,
		},
		{
			filter: o.dependencyFilter(&execution.Execution{InstanceHash: hash.Int(3), TaskKey: "2", ProcessHash: hash.Int(1), StepID: "2"}),
			p:      &p,
			n:      p.Nodes[0],
			res:    false,
			err:    nil,
		},
	}
	for _, test := range tests {
		ok, err := test.filter(test.p, test.n)
		if test.err != nil {
			require.Equal(t, test.err, err)
		} else {
			require.Equal(t, ok, test.res)
		}
	}
}

func TestStart(t *testing.T) {
	event := &mocks.EventSDK{}
	exec := &mocks.ExecutionSDK{}
	process := &mocks.ProcessSDK{}
	o := New(event, exec, process)
	eventListener := &eventsdk.Listener{}
	execListener := &executionsdk.Listener{}
	event.On("GetStream", mock.Anything).Return(eventListener)
	exec.On("GetStream", mock.Anything).Return(execListener)
	err := o.Start()
	eventListener.Close()
	execListener.Close()
	require.NoError(t, err)
}
