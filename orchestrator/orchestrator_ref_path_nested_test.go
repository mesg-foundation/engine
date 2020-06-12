package orchestrator

import (
	"context"
	"testing"
	"time"

	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/stretchr/testify/require"
)

func TestOrchestratorRefPathNested(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	orch, ep, store, _, testInstanceHash, _, _, execChan := newTestOrchestrator(ctx, t)
	defer orch.Stop()

	var (
		testProcessHash hash.Hash
		err             error
	)
	t.Run("create process", func(t *testing.T) {
		testProcessHash, err = store.CreateProcess(
			context.Background(),
			"nested-path-data",
			[]*process.Process_Node{
				{
					Key: "n0",
					Type: &process.Process_Node_Event_{
						Event: &process.Process_Node_Event{
							InstanceHash: testInstanceHash,
							EventKey:     "event_complex_trigger",
						},
					},
				},
				{
					Key: "n1",
					Type: &process.Process_Node_Map_{
						Map: &process.Process_Node_Map{
							Outputs: map[string]*process.Process_Node_Map_Output{
								"msg": {
									Value: &process.Process_Node_Map_Output_Map_{
										Map: &process.Process_Node_Map_Output_Map{
											Outputs: map[string]*process.Process_Node_Map_Output{
												"msg": {
													Value: &process.Process_Node_Map_Output_Ref{
														Ref: &process.Process_Node_Reference{
															NodeKey: "n0",
															Path: &process.Process_Node_Reference_Path{
																Selector: &process.Process_Node_Reference_Path_Key{
																	Key: "msg",
																},
																Path: &process.Process_Node_Reference_Path{
																	Selector: &process.Process_Node_Reference_Path_Key{
																		Key: "msg",
																	},
																},
															},
														},
													},
												},
												"array": {
													Value: &process.Process_Node_Map_Output_List_{
														List: &process.Process_Node_Map_Output_List{
															Outputs: []*process.Process_Node_Map_Output{
																{
																	Value: &process.Process_Node_Map_Output_Ref{
																		Ref: &process.Process_Node_Reference{
																			NodeKey: "n0",
																			Path: &process.Process_Node_Reference_Path{
																				Selector: &process.Process_Node_Reference_Path_Key{
																					Key: "msg",
																				},
																				Path: &process.Process_Node_Reference_Path{
																					Selector: &process.Process_Node_Reference_Path_Key{
																						Key: "array",
																					},
																					Path: &process.Process_Node_Reference_Path{
																						Selector: &process.Process_Node_Reference_Path_Index{
																							Index: 2,
																						},
																					},
																				},
																			},
																		},
																	},
																},
																{
																	Value: &process.Process_Node_Map_Output_Ref{
																		Ref: &process.Process_Node_Reference{
																			NodeKey: "n0",
																			Path: &process.Process_Node_Reference_Path{
																				Selector: &process.Process_Node_Reference_Path_Key{
																					Key: "msg",
																				},
																				Path: &process.Process_Node_Reference_Path{
																					Selector: &process.Process_Node_Reference_Path_Key{
																						Key: "array",
																					},
																					Path: &process.Process_Node_Reference_Path{
																						Selector: &process.Process_Node_Reference_Path_Index{
																							Index: 1,
																						},
																					},
																				},
																			},
																		},
																	},
																},
																{
																	Value: &process.Process_Node_Map_Output_Ref{
																		Ref: &process.Process_Node_Reference{
																			NodeKey: "n0",
																			Path: &process.Process_Node_Reference_Path{
																				Selector: &process.Process_Node_Reference_Path_Key{
																					Key: "msg",
																				},
																				Path: &process.Process_Node_Reference_Path{
																					Selector: &process.Process_Node_Reference_Path_Key{
																						Key: "array",
																					},
																					Path: &process.Process_Node_Reference_Path{
																						Selector: &process.Process_Node_Reference_Path_Index{
																							Index: 0,
																						},
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				{
					Key: "n2",
					Type: &process.Process_Node_Task_{
						Task: &process.Process_Node_Task{
							InstanceHash: testInstanceHash,
							TaskKey:      "task_complex",
						},
					},
				},
				{
					Key: "n3",
					Type: &process.Process_Node_Map_{
						Map: &process.Process_Node_Map{
							Outputs: map[string]*process.Process_Node_Map_Output{
								"msg": {
									Value: &process.Process_Node_Map_Output_Ref{
										Ref: &process.Process_Node_Reference{
											NodeKey: "n2",
											Path: &process.Process_Node_Reference_Path{
												Selector: &process.Process_Node_Reference_Path_Key{
													Key: "msg",
												},
												Path: &process.Process_Node_Reference_Path{
													Selector: &process.Process_Node_Reference_Path_Key{
														Key: "msg",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				{
					Key: "n4",
					Type: &process.Process_Node_Task_{
						Task: &process.Process_Node_Task{
							InstanceHash: testInstanceHash,
							TaskKey:      "task1",
						},
					},
				},
			},
			[]*process.Process_Edge{
				{Src: "n0", Dst: "n1"},
				{Src: "n1", Dst: "n2"},
				{Src: "n2", Dst: "n3"},
				{Src: "n3", Dst: "n4"},
			},
		)
		require.NoError(t, err)
	})
	t.Run("trigger process", func(t *testing.T) {
		_, err := ep.Publish(
			testInstanceHash,
			"event_complex_trigger",
			&types.Struct{
				Fields: map[string]*types.Value{
					"msg": {
						Kind: &types.Value_StructValue{
							StructValue: &types.Struct{
								Fields: map[string]*types.Value{
									"msg": {
										Kind: &types.Value_StringValue{
											StringValue: "complex",
										},
									},
									"timestamp": {
										Kind: &types.Value_NumberValue{
											NumberValue: float64(time.Now().Unix()),
										},
									},
									"array": {
										Kind: &types.Value_ListValue{
											ListValue: &types.ListValue{Values: []*types.Value{
												{Kind: &types.Value_StringValue{StringValue: "first"}},
												{Kind: &types.Value_StringValue{StringValue: "second"}},
												{Kind: &types.Value_StringValue{StringValue: "third"}},
											}},
										},
									},
								},
							},
						},
					},
				},
			},
		)
		require.NoError(t, err)
	})
	t.Run("first ref", func(t *testing.T) {
		var exec *execution.Execution
		t.Run("check created execution", func(t *testing.T) {
			exec = <-execChan
			require.Equal(t, "task_complex", exec.TaskKey)
			require.Equal(t, "n2", exec.NodeKey)
			require.True(t, testProcessHash.Equal(exec.ProcessHash))
			require.Equal(t, execution.Status_InProgress, exec.Status)
			require.Equal(t, "complex", exec.Inputs.Fields["msg"].GetStructValue().Fields["msg"].GetStringValue())
			require.Len(t, exec.Inputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values, 3)
			require.Equal(t, "third", exec.Inputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values[0].GetStringValue())
			require.Equal(t, "second", exec.Inputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values[1].GetStringValue())
			require.Equal(t, "first", exec.Inputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values[2].GetStringValue())
		})
		t.Run("update exec", func(t *testing.T) {
			err := store.UpdateExecution(
				context.Background(),
				exec.Hash,
				time.Now().UnixNano(),
				time.Now().UnixNano(),
				&types.Struct{
					Fields: map[string]*types.Value{
						"msg": {
							Kind: &types.Value_StructValue{
								StructValue: &types.Struct{
									Fields: map[string]*types.Value{
										"msg": {
											Kind: &types.Value_StringValue{
												StringValue: "complex",
											},
										},
										"timestamp": {
											Kind: &types.Value_NumberValue{
												NumberValue: float64(time.Now().Unix()),
											},
										},
										"array": {
											Kind: &types.Value_ListValue{
												ListValue: &types.ListValue{Values: []*types.Value{
													{Kind: &types.Value_StringValue{StringValue: "third"}},
													{Kind: &types.Value_StringValue{StringValue: "second"}},
													{Kind: &types.Value_StringValue{StringValue: "first"}},
												}},
											},
										},
									},
								},
							},
						},
					},
				},
				"",
			)
			require.NoError(t, err)
		})
		t.Run("check completed execution", func(t *testing.T) {
			exec := <-execChan
			require.Equal(t, "task_complex", exec.TaskKey)
			require.Equal(t, "n2", exec.NodeKey)
			require.True(t, testProcessHash.Equal(exec.ProcessHash))
			require.Equal(t, execution.Status_Completed, exec.Status)
			require.Equal(t, "complex", exec.Outputs.Fields["msg"].GetStructValue().Fields["msg"].GetStringValue())
			require.Len(t, exec.Outputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values, 3)
			require.Equal(t, "third", exec.Outputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values[0].GetStringValue())
			require.Equal(t, "second", exec.Outputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values[1].GetStringValue())
			require.Equal(t, "first", exec.Outputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values[2].GetStringValue())
			require.NotEmpty(t, exec.Outputs.Fields["msg"].GetStructValue().Fields["timestamp"].GetNumberValue())
		})
	})
	t.Run("second ref", func(t *testing.T) {
		var exec *execution.Execution
		t.Run("check created execution", func(t *testing.T) {
			exec = <-execChan
			require.Equal(t, "task1", exec.TaskKey)
			require.Equal(t, "n4", exec.NodeKey)
			require.True(t, testProcessHash.Equal(exec.ProcessHash))
			require.Equal(t, execution.Status_InProgress, exec.Status)
			require.Equal(t, "complex", exec.Inputs.Fields["msg"].GetStringValue())
		})
	})
}
