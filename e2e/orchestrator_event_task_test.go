package main

import (
	"context"
	"testing"
	"time"

	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/process"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/stretchr/testify/require"
)

func testOrchestratorEventTask(executionStream pb.Execution_StreamClient, resultStream pb.Result_StreamClient, instanceHash hash.Hash) func(t *testing.T) {
	return func(t *testing.T) {
		var processHash hash.Hash

		t.Run("create process", func(t *testing.T) {
			respProc, err := client.ProcessClient.Create(context.Background(), &pb.CreateProcessRequest{
				Name: "event-task-process",
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
								StringValue: "foo_1",
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
		var execHash hash.Hash
		t.Run("check in progress execution", func(t *testing.T) {
			exec, err := executionStream.Recv()
			require.NoError(t, err)
			require.Equal(t, "task1", exec.TaskKey)
			require.True(t, processHash.Equal(exec.ProcessHash))
			require.Equal(t, "foo_1", exec.Inputs.Fields["msg"].GetStringValue())
			execHash = exec.Hash
		})
		t.Run("check completed execution", func(t *testing.T) {
			res, err := resultStream.Recv()
			require.NoError(t, err)
			require.True(t, res.ExecutionHash.Equal(execHash))
			require.Equal(t, "foo_1", res.GetOutputs().Fields["msg"].GetStringValue())
			require.NotEmpty(t, res.GetOutputs().Fields["timestamp"].GetNumberValue())
		})
		t.Run("delete process", func(t *testing.T) {
			_, err := client.ProcessClient.Delete(context.Background(), &pb.DeleteProcessRequest{Hash: processHash})
			require.NoError(t, err)
		})
	}
}
