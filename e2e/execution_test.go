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
	"google.golang.org/grpc/metadata"
)

func testExecution(t *testing.T) {
	var executionHash hash.Hash
	var execPing *execution.Execution

	ctx := metadata.NewOutgoingContext(context.Background(), passmd)
	t.Run("stream", func(t *testing.T) {
		t.Run("nil filter", func(t *testing.T) {
			_, err := client.ExecutionClient.Stream(context.Background(), &pb.StreamExecutionRequest{})
			require.NoError(t, err)
		})
		t.Run("good", func(t *testing.T) {
			stream, err := client.ExecutionClient.Stream(context.Background(), &pb.StreamExecutionRequest{
				Filter: &pb.StreamExecutionRequest_Filter{
					ExecutorHash: testRunnerHash,
				},
			})
			require.NoError(t, err)
			acknowledgement.WaitForStreamToBeReady(stream)

			resp, err := client.ExecutionClient.Create(ctx, &pb.CreateExecutionRequest{
				TaskKey:      "ping",
				EventHash:    hash.Int(1),
				ExecutorHash: testRunnerHash,
				Inputs: &types.Struct{
					Fields: map[string]*types.Value{
						"msg": {
							Kind: &types.Value_StringValue{
								StringValue: "test",
							},
						},
					},
				},
			})
			require.NoError(t, err)
			executionHash = resp.Hash
			t.Run("receive in progress execution", func(t *testing.T) {
				execPingInProgress, err := stream.Recv()
				require.NoError(t, err)
				require.Equal(t, resp.Hash, execPingInProgress.Hash)
				require.Equal(t, "ping", execPingInProgress.TaskKey)
				require.Equal(t, execution.Status_InProgress, execPingInProgress.Status)
			})
			t.Run("receive completed execution", func(t *testing.T) {
				execPing, err = stream.Recv()
				require.NoError(t, err)
				require.Equal(t, resp.Hash, execPing.Hash)
				require.Equal(t, "ping", execPing.TaskKey)
				require.Equal(t, execution.Status_Completed, execPing.Status)
			})
		})
	})

	t.Run("get", func(t *testing.T) {
		exec, err := client.ExecutionClient.Get(ctx, &pb.GetExecutionRequest{Hash: executionHash})
		require.NoError(t, err)
		require.True(t, exec.Equal(execPing))
	})

	t.Run("list", func(t *testing.T) {
		resp, err := client.ExecutionClient.List(ctx, &pb.ListExecutionRequest{})
		require.NoError(t, err)
		require.Len(t, resp.Executions, 1)
	})
}
