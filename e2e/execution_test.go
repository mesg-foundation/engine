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
	var execPingCompleted *execution.Execution
	ctx := metadata.NewOutgoingContext(context.Background(), passmd)
	t.Run("stream", func(t *testing.T) {
		t.Run("with nil filter", func(t *testing.T) {
			_, err := client.ExecutionClient.Stream(context.Background(), &pb.StreamExecutionRequest{})
			require.NoError(t, err)
		})
		// TODO: no error are returned but it supposed to...
		// t.Run("with not valid filter", func(t *testing.T) {
		// 	t.Run("not found executor", func(t *testing.T) {
		// 		_, err := client.ExecutionClient.Stream(context.Background(), &pb.StreamExecutionRequest{
		// 			Filter: &pb.StreamExecutionRequest_Filter{
		// 				ExecutorHash: hash.Int(1),
		// 			},
		// 		})
		// 		require.EqualError(t, err, "dwdw")
		// 	})
		// 	t.Run("not found instance", func(t *testing.T) {
		// 		_, err := client.ExecutionClient.Stream(context.Background(), &pb.StreamExecutionRequest{
		// 			Filter: &pb.StreamExecutionRequest_Filter{
		// 				InstanceHash: hash.Int(1),
		// 			},
		// 		})
		// 		require.EqualError(t, err, "dwdw")
		// 	})
		// 	t.Run("not found task key", func(t *testing.T) {
		// 		_, err := client.ExecutionClient.Stream(context.Background(), &pb.StreamExecutionRequest{
		// 			Filter: &pb.StreamExecutionRequest_Filter{
		// 				ExecutorHash: testRunnerHash,
		// 				TaskKey:      "do-not-exist",
		// 			},
		// 		})
		// 		require.EqualError(t, err, "service \"test-service\" - task \"do-not-exist\" not found")
		// 	})
		// })
		t.Run("working", func(t *testing.T) {
			stream, err := client.ExecutionClient.Stream(context.Background(), &pb.StreamExecutionRequest{
				Filter: &pb.StreamExecutionRequest_Filter{
					ExecutorHash: testRunnerHash,
				},
			})
			require.NoError(t, err)
			acknowledgement.WaitForStreamToBeReady(stream)
			resp, err := client.ExecutionClient.Create(ctx, &pb.CreateExecutionRequest{
				InstanceHash: testInstanceHash,
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
				execPingCompleted, err = stream.Recv()
				require.NoError(t, err)
				require.Equal(t, resp.Hash, execPingCompleted.Hash)
				require.Equal(t, "ping", execPingCompleted.TaskKey)
				require.Equal(t, execution.Status_Completed, execPingCompleted.Status)
			})
		})
	})

	t.Run("get", func(t *testing.T) {
		execGet, err := client.ExecutionClient.Get(ctx, &pb.GetExecutionRequest{Hash: executionHash})
		require.NoError(t, err)
		require.True(t, execGet.Equal(execPingCompleted))
	})

	// t.Run("list", func(t *testing.T) {
	// 	execList, err := client.ExecutionClient.List(ctx, &pb.GetExecutionRequest{Hash: executionHash})
	// 	require.NoError(t, err)
	// 	require.True(t, execGet.Equal(execPingCompleted))
	// })
}
