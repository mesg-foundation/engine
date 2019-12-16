package main

import (
	"context"
	"testing"

	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/process"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/stretchr/testify/require"
)

func testOrchestratorNestedMap(executionStream pb.Execution_StreamClient, instanceHash hash.Hash) func(t *testing.T) {
	return func(t *testing.T) {
		var (
			processHash hash.Hash
			dataEvent   = &types.Struct{
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
											NumberValue: 101,
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
			}
		)
		t.Run("create process", func(t *testing.T) {
			respProc, err := client.ProcessClient.Create(context.Background(), &pb.CreateProcessRequest{
				Name: "nested-map",
				Nodes: []*process.Process_Node{
					{
						Key: "n0",
						Type: &process.Process_Node_Event_{
							Event: &process.Process_Node_Event{
								InstanceHash: instanceHash,
								EventKey:     "test_event_complex",
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
											Map: &process.Process_Node_Map_Output_Map{Outputs: map[string]*process.Process_Node_Map_Output{
												"msg": {Value: &process.Process_Node_Map_Output_StringConst{
													StringConst: "isAConstant",
												}},
												"array": {Value: &process.Process_Node_Map_Output_List_{
													List: &process.Process_Node_Map_Output_List{Outputs: []*process.Process_Node_Map_Output{
														{Value: &process.Process_Node_Map_Output_StringConst{StringConst: "first-constant"}},
														{Value: &process.Process_Node_Map_Output_StringConst{StringConst: "second-constant"}},
														{Value: &process.Process_Node_Map_Output_StringConst{StringConst: "third-constant"}},
														{Value: &process.Process_Node_Map_Output_StringConst{StringConst: "fourth-constant"}},
													}},
												}},
											}},
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
								InstanceHash: instanceHash,
								TaskKey:      "task_complex",
							},
						},
					},
				},
				Edges: []*process.Process_Edge{
					{Src: "n0", Dst: "n1"},
					{Src: "n1", Dst: "n2"},
				},
			})
			require.NoError(t, err)
			processHash = respProc.Hash
		})
		t.Run("trigger process", func(t *testing.T) {
			_, err := client.EventClient.Create(context.Background(), &pb.CreateEventRequest{
				InstanceHash: instanceHash,
				Key:          "test_event_complex",
				Data:         dataEvent,
			})
			require.NoError(t, err)
		})
		t.Run("first task", func(t *testing.T) {
			t.Run("check in progress execution", func(t *testing.T) {
				exec, err := executionStream.Recv()
				require.NoError(t, err)
				require.Equal(t, "task_complex", exec.TaskKey)
				require.True(t, processHash.Equal(exec.ProcessHash))
				require.Equal(t, execution.Status_InProgress, exec.Status)
				require.Equal(t, "isAConstant", exec.Inputs.Fields["msg"].GetStructValue().Fields["msg"].GetStringValue())
				require.Len(t, exec.Inputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values, 4)
				require.Equal(t, "first-constant", exec.Inputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values[0].GetStringValue())
				require.Equal(t, "second-constant", exec.Inputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values[1].GetStringValue())
				require.Equal(t, "third-constant", exec.Inputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values[2].GetStringValue())
				require.Equal(t, "fourth-constant", exec.Inputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values[3].GetStringValue())
			})
			t.Run("check completed execution", func(t *testing.T) {
				exec, err := executionStream.Recv()
				require.NoError(t, err)
				require.Equal(t, "task_complex", exec.TaskKey)
				require.True(t, processHash.Equal(exec.ProcessHash))
				require.Equal(t, execution.Status_Completed, exec.Status)
				require.Equal(t, "isAConstant", exec.Outputs.Fields["msg"].GetStructValue().Fields["msg"].GetStringValue())
				require.Len(t, exec.Outputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values, 4)
				require.Equal(t, "first-constant", exec.Outputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values[0].GetStringValue())
				require.Equal(t, "second-constant", exec.Outputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values[1].GetStringValue())
				require.Equal(t, "third-constant", exec.Outputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values[2].GetStringValue())
				require.Equal(t, "fourth-constant", exec.Outputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values[3].GetStringValue())
				require.NotEmpty(t, exec.Outputs.Fields["msg"].GetStructValue().Fields["timestamp"].GetNumberValue())
			})
		})
		t.Run("delete process", func(t *testing.T) {
			_, err := client.ProcessClient.Delete(context.Background(), &pb.DeleteProcessRequest{Hash: processHash})
			require.NoError(t, err)
		})
	}
}
