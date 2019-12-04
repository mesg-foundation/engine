package main

import (
	"context"
	"testing"

	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/stretchr/testify/require"
)

func testExecution(t *testing.T) {
	var (
		executionHash hash.Hash
		execPing      *execution.Execution
		taskKey       = "ping"
		eventHash     = hash.Int(1)
		executorHash  = testRunnerHash
		inputs        = &types.Struct{
			Fields: map[string]*types.Value{
				"msg": {
					Kind: &types.Value_StringValue{
						StringValue: "test",
					},
				},
			},
		}
	)

	t.Run("stream", func(t *testing.T) {
		t.Run("nil filter", func(t *testing.T) {
			_, err := client.ExecutionClient.Stream(context.Background(), &pb.StreamExecutionRequest{})
			require.NoError(t, err)
		})
		t.Run("good", func(t *testing.T) {
			stream, err := client.ExecutionClient.Stream(context.Background(), &pb.StreamExecutionRequest{
				Filter: &pb.StreamExecutionRequest_Filter{
					ExecutorHash: executorHash,
				},
			})
			require.NoError(t, err)
			acknowledgement.WaitForStreamToBeReady(stream)

			resp, err := client.ExecutionClient.Create(context.Background(), &pb.CreateExecutionRequest{
				TaskKey:      taskKey,
				EventHash:    eventHash,
				ExecutorHash: executorHash,
				Inputs:       inputs,
			})
			require.NoError(t, err)
			executionHash = resp.Hash
			t.Run("receive in progress execution", func(t *testing.T) {
				execPingInProgress, err := stream.Recv()
				require.NoError(t, err)
				require.Equal(t, resp.Hash, execPingInProgress.Hash)
				require.Equal(t, taskKey, execPingInProgress.TaskKey)
				require.Equal(t, eventHash, execPingInProgress.EventHash)
				require.Equal(t, executorHash, execPingInProgress.ExecutorHash)
				require.Equal(t, execution.Status_InProgress, execPingInProgress.Status)
				require.True(t, inputs.Equal(execPingInProgress.Inputs))
			})
			t.Run("receive completed execution", func(t *testing.T) {
				execPing, err = stream.Recv()
				require.NoError(t, err)
				require.Equal(t, resp.Hash, execPing.Hash)
				require.Equal(t, taskKey, execPing.TaskKey)
				require.Equal(t, eventHash, execPing.EventHash)
				require.Equal(t, executorHash, execPing.ExecutorHash)
				require.Equal(t, execution.Status_Completed, execPing.Status)
				require.True(t, inputs.Equal(execPing.Inputs))
				require.True(t, execPing.Outputs.Equal(&types.Struct{
					Fields: map[string]*types.Value{
						"pong": {
							Kind: &types.Value_StringValue{
								StringValue: "test",
							},
						},
					},
				}))
			})
		})
	})

	t.Run("get", func(t *testing.T) {
		exec, err := client.ExecutionClient.Get(context.Background(), &pb.GetExecutionRequest{Hash: executionHash})
		require.NoError(t, err)
		require.True(t, exec.Equal(execPing))
	})

	t.Run("list", func(t *testing.T) {
		resp, err := client.ExecutionClient.List(context.Background(), &pb.ListExecutionRequest{})
		require.NoError(t, err)
		require.Len(t, resp.Executions, 1)
	})
}
