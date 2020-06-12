package orchestrator

import (
	"context"
	"os"
	"testing"

	"github.com/mesg-foundation/engine/event/publisher"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/instance"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
)

func newTestOrchestrator(ctx context.Context, t *testing.T) (
	orch *Orchestrator,
	ep *publisher.EventPublisher,
	store *storeTest,
	serviceHash hash.Hash,
	instanceHash hash.Hash,
	runnerHash hash.Hash,
	instanceEnvHash hash.Hash,
	execChan <-chan *execution.Execution,
) {
	var err error
	logger := log.NewTMJSONLogger(log.NewSyncWriter(os.Stdout))
	store, err = newStoreTest()
	require.NoError(t, err)
	ep = publisher.New(store)
	orch = New(store, ep, logger)

	go func() {
		require.NoError(t, orch.Start())
	}()

	t.Run("create service", func(t *testing.T) {
		serviceHash, err = createTestService(store)
		require.NoError(t, err)
	})

	t.Run("create runner", func(t *testing.T) {
		instanceEnvHash, err = hash.Random()
		require.NoError(t, err)
		runnerHash, err = store.RegisterRunner(context.Background(), serviceHash, instanceEnvHash)
		require.NoError(t, err)
	})

	t.Run("generate instance hash", func(t *testing.T) {
		inst, err := instance.New(serviceHash, instanceEnvHash)
		require.NoError(t, err)
		instanceHash = inst.Hash
	})

	t.Run("init execution subscriber", func(t *testing.T) {
		execChan, err = store.SubscribeToExecutions(ctx)
		require.NoError(t, err)
	})

	return
}

// XXX: comment test because they were using sdk mocks.
// we don't have sdk now so for now keeping it commented
// TODO: add them later.

// func TestFilter(t *testing.T) {
// 	o := New(&mocks.EventSDK{}, &mocks.ExecutionSDK{}, &mocks.ProcessSDK{}, &mocks.RunnerSDK{})
// 	p := process.Process{
// 		Hash: hash.Int(1),
// 		Nodes: []*process.Process_Node{
// 			{
// 				Key: "1",
// 				Type: &process.Process_Node_Event_{
// 					Event: &process.Process_Node_Event{
// 						InstanceHash: hash.Int(1),
// 						EventKey:     "1",
// 					},
// 				},
// 			},
// 			{
// 				Key: "2",
// 				Type: &process.Process_Node_Task_{Task: &process.Process_Node_Task{
// 					InstanceHash: hash.Int(2),
// 					TaskKey:      "2",
// 				}}},
// 			{
// 				Key: "3",
// 				Type: &process.Process_Node_Task_{Task: &process.Process_Node_Task{
// 					InstanceHash: hash.Int(3),
// 					TaskKey:      "3",
// 				}}},
// 		},
// 		Edges: []*process.Process_Edge{
// 			{Src: "1", Dst: "2"},
// 			{Src: "2", Dst: "3"},
// 		},
// 	}
// 	var tests = []struct {
// 		filter func(wf *process.Process, node *process.Process_Node) (bool, error)
// 		p      *process.Process
// 		n      *process.Process_Node
// 		res    bool
// 		err    error
// 	}{
// 		{
// 			filter: o.eventFilter(&event.Event{InstanceHash: hash.Int(1), Key: "1"}),
// 			n:      p.Nodes[0],
// 			res:    true,
// 			err:    nil,
// 		},
// 		{
// 			filter: o.eventFilter(&event.Event{InstanceHash: hash.Int(1), Key: "2"}),
// 			n:      p.Nodes[0],
// 			res:    false,
// 			err:    nil,
// 		},
// 		{
// 			filter: o.eventFilter(&event.Event{InstanceHash: hash.Int(2), Key: "1"}),
// 			n:      p.Nodes[0],
// 			res:    false,
// 			err:    nil,
// 		},
// 		{
// 			filter: o.eventFilter(&event.Event{InstanceHash: hash.Int(2), Key: "1"}),
// 			n:      p.Nodes[1],
// 			res:    false,
// 			err:    nil,
// 		},
// 		{
// 			filter: o.resultFilter(&execution.Execution{InstanceHash: hash.Int(1), TaskKey: "1"}),
// 			n: &process.Process_Node{Type: &process.Process_Node_Result_{Result: &process.Process_Node_Result{
// 				InstanceHash: hash.Int(1),
// 				TaskKey:      "1",
// 			}}},
// 			res: true,
// 			err: nil,
// 		},
// 		{
// 			filter: o.resultFilter(&execution.Execution{InstanceHash: hash.Int(1), TaskKey: "1"}),
// 			n: &process.Process_Node{Type: &process.Process_Node_Result_{Result: &process.Process_Node_Result{
// 				InstanceHash: hash.Int(1),
// 				TaskKey:      "2",
// 			}}},
// 			res: false,
// 			err: nil,
// 		},
// 		{
// 			filter: o.resultFilter(&execution.Execution{InstanceHash: hash.Int(1), TaskKey: "1"}),
// 			n: &process.Process_Node{Type: &process.Process_Node_Result_{Result: &process.Process_Node_Result{
// 				InstanceHash: hash.Int(2),
// 				TaskKey:      "1",
// 			}}},
// 			res: false,
// 			err: nil,
// 		},
// 		{
// 			filter: o.resultFilter(&execution.Execution{InstanceHash: hash.Int(1), TaskKey: "1"}),
// 			n:      p.Nodes[0],
// 			res:    false,
// 			err:    nil,
// 		},
// 		{
// 			filter: o.dependencyFilter(&execution.Execution{InstanceHash: hash.Int(3), TaskKey: "2", ProcessHash: hash.Int(1), NodeKey: "2"}),
// 			p:      &p,
// 			n:      p.Nodes[2],
// 			res:    true,
// 			err:    nil,
// 		},
// 		{
// 			filter: o.dependencyFilter(&execution.Execution{InstanceHash: hash.Int(3), TaskKey: "2", ProcessHash: hash.Int(2), NodeKey: "2"}),
// 			p:      &p,
// 			n:      p.Nodes[2],
// 			res:    false,
// 			err:    nil,
// 		},
// 		{
// 			filter: o.dependencyFilter(&execution.Execution{InstanceHash: hash.Int(3), TaskKey: "2", ProcessHash: hash.Int(1), NodeKey: "1"}),
// 			p:      &p,
// 			n:      p.Nodes[2],
// 			res:    false,
// 			err:    nil,
// 		},
// 		{
// 			filter: o.dependencyFilter(&execution.Execution{InstanceHash: hash.Int(3), TaskKey: "2", ProcessHash: hash.Int(1), NodeKey: "2"}),
// 			p:      &p,
// 			n:      p.Nodes[0],
// 			res:    false,
// 			err:    nil,
// 		},
// 	}
// 	for _, test := range tests {
// 		ok, err := test.filter(test.p, test.n)
// 		if test.err != nil {
// 			require.Equal(t, test.err, err)
// 		} else {
// 			require.Equal(t, ok, test.res)
// 		}
// 	}
// }

// func TestFindNode(t *testing.T) {
// 	o := New(&mocks.EventSDK{}, &mocks.ExecutionSDK{}, &mocks.ProcessSDK{}, &mocks.RunnerSDK{})
// 	data := &process.Process{
// 		Hash: hash.Int(1),
// 		Nodes: []*process.Process_Node{
// 			{
// 				Key: "1",
// 				Type: &process.Process_Node_Event_{
// 					Event: &process.Process_Node_Event{},
// 				},
// 			},
// 		},
// 	}
// 	require.Len(t, o.findNodes(data, func(p *process.Process, n *process.Process_Node) (bool, error) {
// 		return true, nil
// 	}), 1)
// 	require.Len(t, o.findNodes(data, func(p *process.Process, n *process.Process_Node) (bool, error) {
// 		return n.Key == "1", nil
// 	}), 1)
// 	require.Len(t, o.findNodes(data, func(p *process.Process, n *process.Process_Node) (bool, error) {
// 		return n.Key == "2", nil
// 	}), 0)
// }

// // func TestProcessMap(t *testing.T) {
// // 	e := &mocks.ExecutionSDK{}
// // 	o := New(&mocks.EventSDK{}, e, &mocks.ProcessSDK{})
// // 	exec := &execution.Execution{
// // 		ProcessHash: hash.Int(1),
// // 		NodeKey:      "1",
// // 		ParentHash:  hash.Int(2),
// // 		Outputs: &types.Struct{
// // 			Fields: map[string]*types.Value{
// // 				"xxx": &types.Value{
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
// 	o := New(&mocks.EventSDK{}, e, &mocks.ProcessSDK{}, &mocks.RunnerSDK{})
// 	exec := &execution.Execution{
// 		ProcessHash: hash.Int(2),
// 		NodeKey:     "2",
// 		ParentHash:  hash.Int(3),
// 		Outputs: &types.Struct{
// 			Fields: map[string]*types.Value{
// 				"xxx": {
// 					Kind: &types.Value_StringValue{StringValue: "str"},
// 				},
// 			},
// 		},
// 	}
// 	// Different processes
// 	_, err := o.resolveInput(hash.Int(1), exec, "2", &process.Process_Node_Reference_Path{Selector: &process.Process_Node_Reference_Path_Key{Key: "xxx"}})
// 	require.Error(t, err)
// 	// Different steps, should return the value of the data
// 	val, err := o.resolveInput(hash.Int(2), exec, "2", &process.Process_Node_Reference_Path{Selector: &process.Process_Node_Reference_Path_Key{Key: "xxx"}})
// 	require.NoError(t, err)
// 	require.Equal(t, val, exec.Outputs.Fields["xxx"])
// 	// Invalid execution parent hash
// 	e.On("Get", mock.Anything).Once().Return(nil, fmt.Errorf("err"))
// 	_, err = o.resolveInput(hash.Int(2), exec, "-", &process.Process_Node_Reference_Path{Selector: &process.Process_Node_Reference_Path_Key{Key: "xxx"}})
// 	require.Error(t, err)
// 	// Output from a previous exec
// 	execMock := &execution.Execution{
// 		NodeKey:     "3",
// 		ProcessHash: hash.Int(2),
// 		Outputs: &types.Struct{
// 			Fields: map[string]*types.Value{
// 				"yyy": {
// 					Kind: &types.Value_StringValue{StringValue: "str2"},
// 				},
// 			},
// 		},
// 	}
// 	e.On("Get", mock.Anything).Once().Return(execMock, nil)
// 	val, err = o.resolveInput(hash.Int(2), exec, "3", &process.Process_Node_Reference_Path{Selector: &process.Process_Node_Reference_Path_Key{Key: "yyy"}})
// 	require.NoError(t, err)
// 	require.Equal(t, val, execMock.Outputs.Fields["yyy"])
// }

// func TestProcessTask(t *testing.T) {
// 	e := &mocks.ExecutionSDK{}
// 	e.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(nil, nil)
// 	r := &mocks.RunnerSDK{}
// 	r.On("List", mock.Anything).Once().Return([]*runner.Runner{{Hash: hash.Int(1)}}, nil)
// 	o := New(&mocks.EventSDK{}, e, &mocks.ProcessSDK{}, r)
// 	err := o.processTask("-", &process.Process_Node_Task{
// 		InstanceHash: hash.Int(1),
// 		TaskKey:      "-",
// 	}, &process.Process{
// 		Hash: hash.Int(2),
// 	}, &execution.Execution{
// 		Hash: hash.Int(3),
// 	}, nil, &types.Struct{
// 		Fields: map[string]*types.Value{},
// 	})
// 	require.NoError(t, err)
// 	e.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(nil, fmt.Errorf("error"))
// 	r.On("List", mock.Anything).Once().Return([]*runner.Runner{{Hash: hash.Int(1)}}, nil)
// 	err = o.processTask("-", &process.Process_Node_Task{
// 		InstanceHash: hash.Int(1),
// 		TaskKey:      "-",
// 	}, &process.Process{
// 		Hash: hash.Int(2),
// 	}, nil, &event.Event{
// 		Hash: hash.Int(3),
// 	}, &types.Struct{
// 		Fields: map[string]*types.Value{},
// 	})
// 	require.Error(t, err)
// }
