package main

import (
	"context"
	"testing"
	"time"

	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/protobuf/types"
	processmodule "github.com/mesg-foundation/engine/x/process"
	"github.com/stretchr/testify/require"
)

func testOrchestratorNestedData(instanceHash hash.Hash) func(t *testing.T) {
	return func(t *testing.T) {
		var (
			processHash     hash.Hash
			executionStream pb.Execution_StreamClient
		)

		t.Run("create process", func(t *testing.T) {
			processHash = lcdBroadcastMsg(t, processmodule.MsgCreate{
				Owner: engineAddress,
				Name:  "nested-data",
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
				},
			})
		})
		data := &types.Struct{
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
		}
		t.Run("create execution stream", func(t *testing.T) {
			var err error
			executionStream, err = client.ExecutionClient.Stream(context.Background(), &pb.StreamExecutionRequest{})
			require.NoError(t, err)
			acknowledgement.WaitForStreamToBeReady(executionStream)
		})
		t.Run("trigger process", func(t *testing.T) {
			_, err := client.EventClient.Create(context.Background(), &pb.CreateEventRequest{
				InstanceHash: instanceHash,
				Key:          "test_event_complex",
				Data:         data,
			})
			require.NoError(t, err)
		})
		t.Run("check in progress execution", func(t *testing.T) {
			exec, err := executionStream.Recv()
			require.NoError(t, err)
			require.Equal(t, "task_complex", exec.TaskKey)
			require.Equal(t, "n1", exec.NodeKey)
			require.True(t, processHash.Equal(exec.ProcessHash))
			require.Equal(t, execution.Status_InProgress, exec.Status)
			require.True(t, data.Equal(exec.Inputs))
		})
		t.Run("check completed execution", func(t *testing.T) {
			exec, err := executionStream.Recv()
			require.NoError(t, err)
			require.Equal(t, "task_complex", exec.TaskKey)
			require.Equal(t, "n1", exec.NodeKey)
			require.True(t, processHash.Equal(exec.ProcessHash))
			require.Equal(t, execution.Status_Completed, exec.Status)
			require.True(t, data.Equal(exec.Inputs))
			require.Equal(t, "complex", exec.Outputs.Fields["msg"].GetStructValue().Fields["msg"].GetStringValue())
			require.Len(t, exec.Outputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values, 3)
			require.Equal(t, "first", exec.Outputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values[0].GetStringValue())
			require.Equal(t, "second", exec.Outputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values[1].GetStringValue())
			require.Equal(t, "third", exec.Outputs.Fields["msg"].GetStructValue().Fields["array"].GetListValue().Values[2].GetStringValue())
			require.NotEmpty(t, exec.Outputs.Fields["msg"].GetStructValue().Fields["timestamp"].GetNumberValue())
		})
		t.Run("delete process", func(t *testing.T) {
			lcdBroadcastMsg(t, processmodule.MsgDelete{
				Owner: engineAddress,
				Hash:  processHash,
			})
		})
	}
}
