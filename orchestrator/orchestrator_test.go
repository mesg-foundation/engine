package orchestrator

import (
	"fmt"
	"testing"

	"github.com/mesg-foundation/engine/event"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/orchestrator/mocks"
	"github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/result"
	"github.com/mesg-foundation/engine/runner"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFilter(t *testing.T) {
	o := New(&mocks.EventSDK{}, &mocks.ExecutionSDK{}, &mocks.ResultSDK{}, &mocks.ProcessSDK{}, &mocks.RunnerSDK{}, "", "")
	p := process.Process{
		Hash: hash.Int(1),
		Nodes: []*process.Process_Node{
			{
				Key: "1",
				Type: &process.Process_Node_Event_{
					Event: &process.Process_Node_Event{
						InstanceHash: hash.Int(1),
						EventKey:     "1",
					},
				},
			},
			{
				Key: "2",
				Type: &process.Process_Node_Task_{Task: &process.Process_Node_Task{
					InstanceHash: hash.Int(2),
					TaskKey:      "2",
				}}},
			{
				Key: "3",
				Type: &process.Process_Node_Task_{Task: &process.Process_Node_Task{
					InstanceHash: hash.Int(3),
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
			filter: o.dependencyFilter(&execution.Execution{InstanceHash: hash.Int(3), TaskKey: "2", ProcessHash: hash.Int(1), NodeKey: "2"}),
			p:      &p,
			n:      p.Nodes[2],
			res:    true,
			err:    nil,
		},
		{
			filter: o.dependencyFilter(&execution.Execution{InstanceHash: hash.Int(3), TaskKey: "2", ProcessHash: hash.Int(2), NodeKey: "2"}),
			p:      &p,
			n:      p.Nodes[2],
			res:    false,
			err:    nil,
		},
		{
			filter: o.dependencyFilter(&execution.Execution{InstanceHash: hash.Int(3), TaskKey: "2", ProcessHash: hash.Int(1), NodeKey: "1"}),
			p:      &p,
			n:      p.Nodes[2],
			res:    false,
			err:    nil,
		},
		{
			filter: o.dependencyFilter(&execution.Execution{InstanceHash: hash.Int(3), TaskKey: "2", ProcessHash: hash.Int(1), NodeKey: "2"}),
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

func TestFindNode(t *testing.T) {
	o := New(&mocks.EventSDK{}, &mocks.ExecutionSDK{}, &mocks.ResultSDK{}, &mocks.ProcessSDK{}, &mocks.RunnerSDK{}, "", "")
	data := &process.Process{
		Hash: hash.Int(1),
		Nodes: []*process.Process_Node{
			{
				Key: "1",
				Type: &process.Process_Node_Event_{
					Event: &process.Process_Node_Event{},
				},
			},
		},
	}
	require.Len(t, o.findNodes(data, func(p *process.Process, n *process.Process_Node) (bool, error) {
		return true, nil
	}), 1)
	require.Len(t, o.findNodes(data, func(p *process.Process, n *process.Process_Node) (bool, error) {
		return n.Key == "1", nil
	}), 1)
	require.Len(t, o.findNodes(data, func(p *process.Process, n *process.Process_Node) (bool, error) {
		return n.Key == "2", nil
	}), 0)
}

// // func TestProcessMap(t *testing.T) {
// // 	e := &mocks.ExecutionSDK{}
// // 	o := New(&mocks.EventSDK{}, e, &mocks.ProcessSDK{})
// // 	exec := &execution.Execution{
// // 		ProcessHash: hash.Int(1),
// // 		NodeKey:      "1",
// // 		ParentHash:  hash.Int(2),
// // 		Outputs: &types.Struct{
// // 			Fields: map[string]*types.Value{
// // 				"outputKeyXXX": &types.Value{
// // 					Kind: &types.Value_StringValue{StringValue: "str"},
// // 				},
// // 			},
// // 		},
// // 	}
// // 	o.processMap(&process.Process_Node_Map{
// // 		Outputs: []*process.Process_Node_Map_Output{},
// // 	}, )
// // }

// func TestResolveInput(t *testing.T) {
// 	e := &mocks.ExecutionSDK{}
// 	o := New(&mocks.EventSDK{}, e, &mocks.ProcessSDK{}, &mocks.RunnerSDK{}, "", "")
// 	exec := &execution.Execution{
// 		ProcessHash: hash.Int(1),
// 		NodeKey:      "1",
// 		ParentHash:  hash.Int(2),
// 		Outputs: &types.Struct{
// 			Fields: map[string]*types.Value{
// 				"outputKeyXXX": {
// 					Kind: &types.Value_StringValue{StringValue: "str"},
// 				},
// 			},
// 		},
// 	}
// 	// Different processes
// 	_, err := o.resolveInput(hash.Int(1), exec, "2", "outputKeyXXX")
// 	require.Error(t, err)
// 	// Different steps, should return the value of the data
// 	val, err := o.resolveInput(hash.Int(2), exec, "2", "outputKeyXXX")
// 	require.NoError(t, err)
// 	require.Equal(t, val, exec.Outputs.Fields["outputKeyXXX"])
// 	// Invalid execution parent hash
// 	e.On("Get", mock.Anything).Once().Return(nil, fmt.Errorf("err"))
// 	_, err = o.resolveInput(hash.Int(2), exec, "-", "outputKeyXXX")
// 	require.Error(t, err)
// 	// Output from a previous exec
// 	execMock := &execution.Execution{
// 		NodeKey:      "nodeKey3",
// 		ProcessHash: hash.Int(2),
// 		Outputs: &types.Struct{
// 			Fields: map[string]*types.Value{
// 				"outputKeyYYY": {
// 					Kind: &types.Value_StringValue{StringValue: "str2"},
// 				},
// 			},
// 		},
// 	}
// 	e.On("Get", mock.Anything).Once().Return(execMock, nil)
// 	val, err = o.resolveInput(hash.Int(2), exec, "nodeKey3", "outputKeyYYY")
// 	require.NoError(t, err)
// 	require.Equal(t, val, execMock.Outputs.Fields["outputKeyYYY"])
// }

func TestResolveInput(t *testing.T) {
	var (
		e                = &mocks.ExecutionSDK{}
		r                = &mocks.ResultSDK{}
		o                = New(&mocks.EventSDK{}, e, r, &mocks.ProcessSDK{}, &mocks.RunnerSDK{}, "", "")
		processHash      = hash.Int(2)
		instanceHash     = hash.Int(1)
		parentResultHash = hash.Int(3)
		exec             = execution.New(
			processHash,
			instanceHash,
			parentResultHash,
			nil,
			"nodeKey2",
			"taskKey",
			nil,
			nil,
			nil,
		)
		res = result.NewWithOutputs(
			exec.Hash,
			&types.Struct{
				Fields: map[string]*types.Value{
					"outputKeyXXX": {
						Kind: &types.Value_StringValue{StringValue: "str"},
					},
				},
			},
		)
	)
	// Different processes
	_, err := o.resolveInput(instanceHash, exec, res, "nodeKey2", "outputKeyoutputKeyXXX")
	require.Error(t, err)
	// Different steps, should return the value of the data
	val, err := o.resolveInput(processHash, exec, res, "nodeKey2", "outputKeyXXX")
	require.NoError(t, err)
	require.Equal(t, val, res.GetOutputs().Fields["outputKeyXXX"])
	// Invalid execution parent hash
	r.On("Get", mock.Anything).Once().Return(nil, fmt.Errorf("err"))
	_, err = o.resolveInput(processHash, exec, res, "-", "outputKeyXXX")
	require.Error(t, err)
	// Output from a previous exec
	var (
		execMock = execution.New(
			processHash,
			instanceHash,
			parentResultHash,
			nil,
			"nodeKey3",
			"taskKey",
			nil,
			nil,
			nil,
		)
		resMock = result.NewWithOutputs(
			execMock.Hash,
			&types.Struct{
				Fields: map[string]*types.Value{
					"outputKeyYYY": {
						Kind: &types.Value_StringValue{StringValue: "str2"},
					},
				},
			},
		)
	)
	e.On("Get", mock.Anything).Once().Return(execMock, nil)
	r.On("Get", mock.Anything).Once().Return(resMock, nil)
	val, err = o.resolveInput(processHash, exec, res, "nodeKey3", "outputKeyYYY")
	require.NoError(t, err)
	require.Equal(t, val, resMock.GetOutputs().Fields["outputKeyYYY"])
}

func TestProcessTask(t *testing.T) {
	e := &mocks.ExecutionSDK{}
	e.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(nil, nil)
	r := &mocks.RunnerSDK{}
	r.On("List", mock.Anything).Once().Return([]*runner.Runner{{Hash: hash.Int(1)}}, nil)
	o := New(&mocks.EventSDK{}, e, &mocks.ResultSDK{}, &mocks.ProcessSDK{}, r, "", "")
	err := o.processTask("-", &process.Process_Node_Task{
		InstanceHash: hash.Int(1),
		TaskKey:      "-",
	}, &process.Process{
		Hash: hash.Int(2),
	}, &result.Result{
		Hash: hash.Int(3),
	}, nil, &types.Struct{
		Fields: map[string]*types.Value{},
	})
	require.NoError(t, err)
	e.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(nil, fmt.Errorf("error"))
	r.On("List", mock.Anything).Once().Return([]*runner.Runner{{Hash: hash.Int(1)}}, nil)
	err = o.processTask("-", &process.Process_Node_Task{
		InstanceHash: hash.Int(1),
		TaskKey:      "-",
	}, &process.Process{
		Hash: hash.Int(2),
	}, nil, &event.Event{
		Hash: hash.Int(3),
	}, &types.Struct{
		Fields: map[string]*types.Value{},
	})
	require.Error(t, err)
}
