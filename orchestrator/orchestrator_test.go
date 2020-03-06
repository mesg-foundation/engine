package orchestrator

// XXX: comment test because they were using sdk mocks.
// we don't have sdk now so for now keeping it commented
// TODO: add them later.

// func TestFilter(t *testing.T) {
// 	o := New(&mocks.EventSDK{}, &mocks.ExecutionSDK{}, &mocks.ProcessSDK{}, &mocks.RunnerSDK{})
// 	p := process.Process{
// 		Hash: address.ProcAddress(crypto.AddressHash([]byte("1"))),
// 		Nodes: []*process.Process_Node{
// 			{
// 				Key: "1",
// 				Type: &process.Process_Node_Event_{
// 					Event: &process.Process_Node_Event{
// 						InstanceHash: address.InstAddress(crypto.AddressHash([]byte("1"))),
// 						EventKey:     "1",
// 					},
// 				},
// 			},
// 			{
// 				Key: "2",
// 				Type: &process.Process_Node_Task_{Task: &process.Process_Node_Task{
// 					InstanceHash: address.InstAddress(crypto.AddressHash([]byte("2"))),
// 					TaskKey:      "2",
// 				}}},
// 			{
// 				Key: "3",
// 				Type: &process.Process_Node_Task_{Task: &process.Process_Node_Task{
// 					InstanceHash: address.InstAddress(crypto.AddressHash([]byte("3"))),
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
// 			filter: o.eventFilter(&event.Event{InstanceHash: address.InstAddress(crypto.AddressHash([]byte("1"))), Key: "1"}),
// 			n:      p.Nodes[0],
// 			res:    true,
// 			err:    nil,
// 		},
// 		{
// 			filter: o.eventFilter(&event.Event{InstanceHash: address.InstAddress(crypto.AddressHash([]byte("1"))), Key: "2"}),
// 			n:      p.Nodes[0],
// 			res:    false,
// 			err:    nil,
// 		},
// 		{
// 			filter: o.eventFilter(&event.Event{InstanceHash: address.InstAddress(crypto.AddressHash([]byte("2"))), Key: "1"}),
// 			n:      p.Nodes[0],
// 			res:    false,
// 			err:    nil,
// 		},
// 		{
// 			filter: o.eventFilter(&event.Event{InstanceHash: address.InstAddress(crypto.AddressHash([]byte("2"))), Key: "1"}),
// 			n:      p.Nodes[1],
// 			res:    false,
// 			err:    nil,
// 		},
// 		{
// 			filter: o.resultFilter(&execution.Execution{InstanceHash: address.InstAddress(crypto.AddressHash([]byte("1"))), TaskKey: "1"}),
// 			n: &process.Process_Node{Type: &process.Process_Node_Result_{Result: &process.Process_Node_Result{
// 				InstanceHash: address.InstAddress(crypto.AddressHash([]byte("1"))),
// 				TaskKey:      "1",
// 			}}},
// 			res: true,
// 			err: nil,
// 		},
// 		{
// 			filter: o.resultFilter(&execution.Execution{InstanceHash: address.InstAddress(crypto.AddressHash([]byte("1"))), TaskKey: "1"}),
// 			n: &process.Process_Node{Type: &process.Process_Node_Result_{Result: &process.Process_Node_Result{
// 				InstanceHash: address.InstAddress(crypto.AddressHash([]byte("1"))),
// 				TaskKey:      "2",
// 			}}},
// 			res: false,
// 			err: nil,
// 		},
// 		{
// 			filter: o.resultFilter(&execution.Execution{InstanceHash: address.InstAddress(crypto.AddressHash([]byte("1"))), TaskKey: "1"}),
// 			n: &process.Process_Node{Type: &process.Process_Node_Result_{Result: &process.Process_Node_Result{
// 				InstanceHash: address.InstAddress(crypto.AddressHash([]byte("2"))),
// 				TaskKey:      "1",
// 			}}},
// 			res: false,
// 			err: nil,
// 		},
// 		{
// 			filter: o.resultFilter(&execution.Execution{InstanceHash: address.InstAddress(crypto.AddressHash([]byte("1"))), TaskKey: "1"}),
// 			n:      p.Nodes[0],
// 			res:    false,
// 			err:    nil,
// 		},
// 		{
// 			filter: o.dependencyFilter(&execution.Execution{InstanceHash: address.InstAddress(crypto.AddressHash([]byte("3"))), TaskKey: "2", ProcessHash: address.InstAddress(crypto.AddressHash([]byte("1"))), NodeKey: "2"}),
// 			p:      &p,
// 			n:      p.Nodes[2],
// 			res:    true,
// 			err:    nil,
// 		},
// 		{
// 			filter: o.dependencyFilter(&execution.Execution{InstanceHash: address.InstAddress(crypto.AddressHash([]byte("3"))), TaskKey: "2", ProcessHash: address.InstAddress(crypto.AddressHash([]byte("2"))), NodeKey: "2"}),
// 			p:      &p,
// 			n:      p.Nodes[2],
// 			res:    false,
// 			err:    nil,
// 		},
// 		{
// 			filter: o.dependencyFilter(&execution.Execution{InstanceHash: address.InstAddress(crypto.AddressHash([]byte("3"))), TaskKey: "2", ProcessHash: address.InstAddress(crypto.AddressHash([]byte("1"))), NodeKey: "1"}),
// 			p:      &p,
// 			n:      p.Nodes[2],
// 			res:    false,
// 			err:    nil,
// 		},
// 		{
// 			filter: o.dependencyFilter(&execution.Execution{InstanceHash: address.InstAddress(crypto.AddressHash([]byte("3"))), TaskKey: "2", ProcessHash: address.InstAddress(crypto.AddressHash([]byte("1"))), NodeKey: "2"}),
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
// 		Hash: address.ProcAddress(crypto.AddressHash([]byte("1"))),
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
// // 		ProcessHash: address.ProcAddress(crypto.AddressHash([]byte("1"))),
// // 		NodeKey:      "1",
// // 		ParentHash:  address.ExecAddress(crypto.AddressHash([]byte("2"))),
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
// 		ProcessHash: address.ProcAddress(crypto.AddressHash([]byte("2"))),
// 		NodeKey:     "2",
// 		ParentHash:  address.ExecAddress(crypto.AddressHash([]byte("3"))),
// 		Outputs: &types.Struct{
// 			Fields: map[string]*types.Value{
// 				"xxx": {
// 					Kind: &types.Value_StringValue{StringValue: "str"},
// 				},
// 			},
// 		},
// 	}
// 	// Different processes
// 	_, err := o.resolveInput(address.ProcAddress(crypto.AddressHash([]byte("1"))), exec, "2", &process.Process_Node_Map_Output_Reference_Path{Selector: &process.Process_Node_Map_Output_Reference_Path_Key{Key: "xxx"}})
// 	require.Error(t, err)
// 	// Different steps, should return the value of the data
// 	val, err := o.resolveInput(address.ProcAddress(crypto.AddressHash([]byte("2"))), exec, "2", &process.Process_Node_Map_Output_Reference_Path{Selector: &process.Process_Node_Map_Output_Reference_Path_Key{Key: "xxx"}})
// 	require.NoError(t, err)
// 	require.Equal(t, val, exec.Outputs.Fields["xxx"])
// 	// Invalid execution parent hash
// 	e.On("Get", mock.Anything).Once().Return(nil, fmt.Errorf("err"))
// 	_, err = o.resolveInput(address.ProcAddress(crypto.AddressHash([]byte("2"))), exec, "-", &process.Process_Node_Map_Output_Reference_Path{Selector: &process.Process_Node_Map_Output_Reference_Path_Key{Key: "xxx"}})
// 	require.Error(t, err)
// 	// Output from a previous exec
// 	execMock := &execution.Execution{
// 		NodeKey:     "3",
// 		ProcessHash: address.ProcAddress(crypto.AddressHash([]byte("2"))),
// 		Outputs: &types.Struct{
// 			Fields: map[string]*types.Value{
// 				"yyy": {
// 					Kind: &types.Value_StringValue{StringValue: "str2"},
// 				},
// 			},
// 		},
// 	}
// 	e.On("Get", mock.Anything).Once().Return(execMock, nil)
// 	val, err = o.resolveInput(address.ProcAddress(crypto.AddressHash([]byte("2"))), exec, "3", &process.Process_Node_Map_Output_Reference_Path{Selector: &process.Process_Node_Map_Output_Reference_Path_Key{Key: "yyy"}})
// 	require.NoError(t, err)
// 	require.Equal(t, val, execMock.Outputs.Fields["yyy"])
// }

// func TestProcessTask(t *testing.T) {
// 	e := &mocks.ExecutionSDK{}
// 	e.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(nil, nil)
// 	r := &mocks.RunnerSDK{}
// 	r.On("List", mock.Anything).Once().Return([]*runner.Runner{{Hash: address.RunAddress(crypto.AddressHash([]byte("1")))}}, nil)
// 	o := New(&mocks.EventSDK{}, e, &mocks.ProcessSDK{}, r)
// 	err := o.processTask("-", &process.Process_Node_Task{
// 		InstanceHash: address.InstAddress(crypto.AddressHash([]byte("1"))),
// 		TaskKey:      "-",
// 	}, &process.Process{
// 		Hash: address.ProcAddress(crypto.AddressHash([]byte("2"))),
// 	}, &execution.Execution{
// 		Hash: address.ProcAddress(crypto.AddressHash([]byte("3"))),
// 	}, nil, &types.Struct{
// 		Fields: map[string]*types.Value{},
// 	})
// 	require.NoError(t, err)
// 	e.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(nil, fmt.Errorf("error"))
// 	r.On("List", mock.Anything).Once().Return([]*runner.Runner{{Hash: address.RunAddress(crypto.AddressHash([]byte("1")))}}, nil)
// 	err = o.processTask("-", &process.Process_Node_Task{
// 		InstanceHash: address.InstAddress(crypto.AddressHash([]byte("1"))),
// 		TaskKey:      "-",
// 	}, &process.Process{
// 		Hash: address.ProcAddress(crypto.AddressHash([]byte("2"))),
// 	}, nil, &event.Event{
// 		Hash: address.ProcAddress(crypto.AddressHash([]byte("3"))),
// 	}, &types.Struct{
// 		Fields: map[string]*types.Value{},
// 	})
// 	require.Error(t, err)
// }
