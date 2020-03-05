package main

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/process"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
)

func testOrchestratorResultTask(executionStream pb.Execution_StreamClient, runnerHash sdk.AccAddress, instanceHash sdk.AccAddress) func(t *testing.T) {
	return func(t *testing.T) {
		var processHash sdk.AccAddress

		t.Run("create process", func(t *testing.T) {
			respProc, err := client.ProcessClient.Create(context.Background(), &pb.CreateProcessRequest{
				Name: "result-task-process",
				Nodes: []*process.Process_Node{
					{
						Key: "n0",
						Type: &process.Process_Node_Result_{
							Result: &process.Process_Node_Result{
								InstanceHash: instanceHash,
								TaskKey:      "task1",
							},
						},
					},
					{
						Key: "n1",
						Type: &process.Process_Node_Task_{
							Task: &process.Process_Node_Task{
								InstanceHash: instanceHash,
								TaskKey:      "task2",
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
			_, err := client.ExecutionClient.Create(context.Background(), &pb.CreateExecutionRequest{
				TaskKey:      "task1",
				EventHash:    sdk.AccAddress(crypto.AddressHash([]byte("11010101011"))),
				ExecutorHash: runnerHash,
				Inputs: &types.Struct{
					Fields: map[string]*types.Value{
						"msg": {
							Kind: &types.Value_StringValue{
								StringValue: "foo_2",
							},
						},
					},
				},
			})
			require.NoError(t, err)
		})
		t.Run("check trigger process execution", func(t *testing.T) {
			t.Run("in progress", func(t *testing.T) {
				exec, err := executionStream.Recv()
				require.NoError(t, err)
				require.Equal(t, "task1", exec.TaskKey)
				require.Equal(t, "", exec.NodeKey)
				require.True(t, sdk.AccAddress(crypto.AddressHash([]byte("11010101011"))).Equals(exec.EventHash))
				require.Equal(t, execution.Status_InProgress, exec.Status)
				require.True(t, exec.Inputs.Equal(&types.Struct{
					Fields: map[string]*types.Value{
						"msg": {
							Kind: &types.Value_StringValue{
								StringValue: "foo_2",
							},
						},
					},
				}))
			})
			t.Run("completed", func(t *testing.T) {
				exec, err := executionStream.Recv()
				require.NoError(t, err)
				require.Equal(t, "task1", exec.TaskKey)
				require.Equal(t, "", exec.NodeKey)
				require.True(t, sdk.AccAddress(crypto.AddressHash([]byte("11010101011"))).Equals(exec.EventHash))
				require.Equal(t, execution.Status_Completed, exec.Status)
				require.Equal(t, "foo_2", exec.Outputs.Fields["msg"].GetStringValue())
				require.NotEmpty(t, exec.Outputs.Fields["timestamp"].GetNumberValue())
			})
		})
		t.Run("check in progress execution", func(t *testing.T) {
			exec, err := executionStream.Recv()
			require.NoError(t, err)
			require.Equal(t, "task2", exec.TaskKey)
			require.Equal(t, "n1", exec.NodeKey)
			require.Equal(t, processHash, exec.ProcessHash)
			require.Equal(t, execution.Status_InProgress, exec.Status)
			require.Equal(t, "foo_2", exec.Inputs.Fields["msg"].GetStringValue())
		})
		t.Run("check completed execution", func(t *testing.T) {
			exec, err := executionStream.Recv()
			require.NoError(t, err)
			require.Equal(t, "task2", exec.TaskKey)
			require.Equal(t, "n1", exec.NodeKey)
			require.Equal(t, processHash, exec.ProcessHash)
			require.Equal(t, execution.Status_Completed, exec.Status)
			require.Equal(t, "foo_2", exec.Outputs.Fields["msg"].GetStringValue())
			require.NotEmpty(t, exec.Outputs.Fields["timestamp"].GetNumberValue())
		})
		t.Run("delete process", func(t *testing.T) {
			_, err := client.ProcessClient.Delete(context.Background(), &pb.DeleteProcessRequest{Hash: processHash})
			require.NoError(t, err)
		})
	}
}
