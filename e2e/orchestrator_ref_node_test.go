package main

import (
	"context"
	"testing"
	"time"

	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/process"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/stretchr/testify/require"
)

func testOrchestratorRefNode(executionStream pb.Execution_StreamClient, instanceHash hash.Hash) func(t *testing.T) {
	return func(t *testing.T) {
		var processHash hash.Hash

		t.Run("create process", func(t *testing.T) {
			respProc, err := client.ProcessClient.Create(context.Background(), &pb.CreateProcessRequest{
				Name: "ref-node",
				Nodes: []*process.Process_Node{
					{
						Key: "n0",
						Type: &process.Process_Node_Event_{
							Event: &process.Process_Node_Event{
								InstanceHash: instanceHash,
								EventKey:     "test_event",
							},
						},
					},
					{
						Key: "n1",
						Type: &process.Process_Node_Map_{
							Map: &process.Process_Node_Map{
								Outputs: map[string]*process.Process_Node_Map_Output{
									"msg": {
										Value: &process.Process_Node_Map_Output_Ref{
											Ref: &process.Process_Node_Map_Output_Reference{
												NodeKey: "n0",
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
		var (
			timestamp = float64(time.Now().Unix())
		)
		t.Run("trigger process", func(t *testing.T) {
			_, err := client.EventClient.Create(context.Background(), &pb.CreateEventRequest{
				InstanceHash: instanceHash,
				Key:          "test_event",
				Data: &types.Struct{
					Fields: map[string]*types.Value{
						"msg": {
							Kind: &types.Value_StringValue{
								StringValue: "whatever",
							},
						},
						"timestamp": {
							Kind: &types.Value_NumberValue{
								NumberValue: timestamp,
							},
						},
					},
				},
			})
			require.NoError(t, err)
		})
		t.Run("check in progress execution", func(t *testing.T) {
			exec, err := executionStream.Recv()
			require.NoError(t, err)
			require.Equal(t, "task_complex", exec.TaskKey)
			require.Equal(t, "n2", exec.NodeKey)
			require.True(t, processHash.Equal(exec.ProcessHash))
			require.Equal(t, execution.Status_InProgress, exec.Status)
			require.Equal(t, "whatever", exec.Inputs.Fields["msg"].GetStructValue().Fields["msg"].GetStringValue())
			require.Equal(t, timestamp, exec.Inputs.Fields["msg"].GetStructValue().Fields["timestamp"].GetNumberValue())
		})
		t.Run("check completed execution", func(t *testing.T) {
			exec, err := executionStream.Recv()
			require.NoError(t, err)
			require.Equal(t, "task_complex", exec.TaskKey)
			require.Equal(t, "n2", exec.NodeKey)
			require.True(t, processHash.Equal(exec.ProcessHash))
			require.Equal(t, execution.Status_Completed, exec.Status)
			require.Equal(t, "whatever", exec.Outputs.Fields["msg"].GetStructValue().Fields["msg"].GetStringValue())
			require.Nil(t, exec.Outputs.Fields["msg"].GetStructValue().Fields["array"])
			require.NotEmpty(t, exec.Outputs.Fields["msg"].GetStructValue().Fields["timestamp"].GetNumberValue())
		})
		t.Run("delete process", func(t *testing.T) {
			_, err := client.ProcessClient.Delete(context.Background(), &pb.DeleteProcessRequest{Hash: processHash})
			require.NoError(t, err)
		})
	}
}
