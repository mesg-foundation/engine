package main

import (
	"context"
	"testing"

	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func testExecution(t *testing.T) {
	stream, err := client.ExecutionClient.Stream(context.Background(), &pb.StreamExecutionRequest{
		Filter: &pb.StreamExecutionRequest_Filter{
			InstanceHash: testInstanceHash,
		},
	})
	require.NoError(t, err)
	acknowledgement.WaitForStreamToBeReady(stream)

	ctx := metadata.NewOutgoingContext(context.Background(), passmd)
	resp, err := client.ExecutionClient.Create(ctx, &pb.CreateExecutionRequest{
		InstanceHash: testInstanceHash,
		TaskKey:      "ping",
		Inputs: []*types.Value{
			{
				Kind: &types.Value_StringValue{
					StringValue: "test",
				},
			},
		},
	})
	require.NoError(t, err)

	// recive in progress status
	exec, err := stream.Recv()
	require.NoError(t, err)
	require.Equal(t, resp.Hash, exec.Hash)
	require.Equal(t, "ping", exec.TaskKey)
	require.Equal(t, execution.Status_InProgress, exec.Status)

	// recive completed status
	exec, err = stream.Recv()
	require.NoError(t, err)
	require.Equal(t, resp.Hash, exec.Hash)
	require.Equal(t, "ping", exec.TaskKey)
	require.Equal(t, execution.Status_Completed, exec.Status)

	exec, err = client.ExecutionClient.Get(ctx, &pb.GetExecutionRequest{Hash: resp.Hash})
	require.NoError(t, err)
	require.Equal(t, resp.Hash, exec.Hash)
	require.Equal(t, "ping", exec.TaskKey)
	require.Equal(t, execution.Status_Completed, exec.Status)
}
