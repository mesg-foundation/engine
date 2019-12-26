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

func testOrchestratorMapConst(executionStream pb.Execution_StreamClient, instanceHash hash.Hash) func(t *testing.T) {
	return func(t *testing.T) {
		var processHash hash.Hash

		t.Run("create process", func(t *testing.T) {
			respProc, err := client.ProcessClient.Create(context.Background(), &pb.CreateProcessRequest{
				Name: "map-const",
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
										Value: &process.Process_Node_Map_Output_StringConst{
											StringConst: "itsAConstant",
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
								TaskKey:      "task1",
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
								NumberValue: float64(time.Now().Unix()),
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
			require.Equal(t, "task1", exec.TaskKey)
			require.Equal(t, "n2", exec.NodeKey)
			require.True(t, processHash.Equal(exec.ProcessHash))
			require.Equal(t, execution.Status_InProgress, exec.Status)
			require.Equal(t, "itsAConstant", exec.Inputs.Fields["msg"].GetStringValue())
		})
		t.Run("check completed execution", func(t *testing.T) {
			exec, err := executionStream.Recv()
			require.NoError(t, err)
			require.Equal(t, "task1", exec.TaskKey)
			require.Equal(t, "n2", exec.NodeKey)
			require.True(t, processHash.Equal(exec.ProcessHash))
			require.Equal(t, execution.Status_Completed, exec.Status)
			require.Equal(t, "itsAConstant", exec.Outputs.Fields["msg"].GetStringValue())
			require.NotEmpty(t, exec.Outputs.Fields["timestamp"].GetNumberValue())
		})
		t.Run("delete process", func(t *testing.T) {
			_, err := client.ProcessClient.Delete(context.Background(), &pb.DeleteProcessRequest{Hash: processHash})
			require.NoError(t, err)
		})
	}
}
