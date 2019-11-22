package main

import (
	"context"
	"testing"

	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/stretchr/testify/require"
)

func testExecution(t *testing.T) {
	stream, err := client.ExecutionClient.Stream(context.Background(), &pb.StreamExecutionRequest{
		Filter: &pb.StreamExecutionRequest_Filter{
			InstanceHash: testInstanceHash,
		},
	})
	require.NoError(t, err)
	acknowledgement.WaitForStreamToBeReady(stream)

	resp, err := client.ExecutionClient.Create(context.Background(), &pb.CreateExecutionRequest{
		InstanceHash: testInstanceHash,
		TaskKey:      "ping",
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

	// receive in progress status
	exec, err := stream.Recv()
	require.NoError(t, err)
	require.Equal(t, resp.Hash, exec.Hash)
	require.Equal(t, "ping", exec.TaskKey)
	require.Equal(t, execution.Status_InProgress, exec.Status)

	// receive completed status
	exec, err = stream.Recv()
	require.NoError(t, err)
	require.Equal(t, resp.Hash, exec.Hash)
	require.Equal(t, "ping", exec.TaskKey)
	require.Equal(t, execution.Status_Completed, exec.Status)

	exec, err = client.ExecutionClient.Get(context.Background(), &pb.GetExecutionRequest{Hash: resp.Hash})
	require.NoError(t, err)
	require.Equal(t, resp.Hash, exec.Hash)
	require.Equal(t, "ping", exec.TaskKey)
	require.Equal(t, execution.Status_Completed, exec.Status)
}
