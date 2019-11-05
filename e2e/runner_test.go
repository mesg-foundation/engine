package main

import (
	"context"
	"testing"

	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

var testRunnerHash hash.Hash

func testRunner(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		stream, err := client.EventClient.Stream(context.Background(), &pb.StreamEventRequest{
			Filter: &pb.StreamEventRequest_Filter{
				Key: "test_service_ready",
			},
		})
		require.NoError(t, err)
		acknowledgement.WaitForStreamToBeReady(stream)

		ctx := metadata.NewOutgoingContext(context.Background(), passmd)
		resp, err := client.RunnerClient.Create(ctx, &pb.CreateRunnerRequest{
			ServiceHash: testServiceHash,
			Env:         []string{"BAR=3", "REQUIRED=4"},
		})
		require.NoError(t, err)
		testRunnerHash = resp.Hash

		// wait for service to be ready
		_, err = stream.Recv()
		require.NoError(t, err)
	})

	t.Run("get", func(t *testing.T) {
		resp, err := client.RunnerClient.Get(context.Background(), &pb.GetRunnerRequest{Hash: testRunnerHash})
		require.NoError(t, err)
		require.Equal(t, testRunnerHash, resp.Hash)
		testInstanceHash = resp.InstanceHash
	})

	t.Run("list", func(t *testing.T) {
		resp, err := client.RunnerClient.List(context.Background(), &pb.ListRunnerRequest{InstanceHash: testInstanceHash})
		require.NoError(t, err)
		require.Len(t, resp.Runners, 1)
		require.Equal(t, testInstanceHash, resp.Runners[0].InstanceHash)
		require.Equal(t, testRunnerHash, resp.Runners[0].Hash)
	})
}

func testDeleteRunner(t *testing.T) {
	ctx := metadata.NewOutgoingContext(context.Background(), passmd)
	_, err := client.RunnerClient.Delete(ctx, &pb.DeleteRunnerRequest{Hash: testRunnerHash})
	require.NoError(t, err)

	resp, err := client.RunnerClient.List(context.Background(), &pb.ListRunnerRequest{InstanceHash: testInstanceHash})
	require.NoError(t, err)
	require.Len(t, resp.Runners, 0)
}
