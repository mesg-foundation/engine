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

var testInstanceHash hash.Hash

func testInstance(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		stream, err := client.EventClient.Stream(context.Background(), &pb.StreamEventRequest{
			Filter: &pb.StreamEventRequest_Filter{
				Key: "test_service_ready",
			},
		})
		require.NoError(t, err)
		acknowledgement.WaitForStreamToBeReady(stream)

		ctx := metadata.NewOutgoingContext(context.Background(), passmd)
		resp, err := client.InstanceClient.Create(ctx, &pb.CreateInstanceRequest{
			ServiceHash: testServiceHash,
			Env:         []string{"BAR=3", "REQUIRED=4"},
		})
		require.NoError(t, err)
		testInstanceHash = resp.Hash

		// wait for service to be ready
		_, err = stream.Recv()
		require.NoError(t, err)
	})

	t.Run("get", func(t *testing.T) {
		resp, err := client.InstanceClient.Get(context.Background(), &pb.GetInstanceRequest{Hash: testInstanceHash})
		require.NoError(t, err)
		require.Equal(t, testInstanceHash, resp.Hash)
		require.Equal(t, testServiceHash, resp.ServiceHash)
		require.Equal(t, hash.Dump([]string{"BAR=2", "FOO=1", "REQUIRED", "BAR=3", "REQUIRED=4"}), resp.EnvHash)
	})

	t.Run("list", func(t *testing.T) {
		resp, err := client.InstanceClient.List(context.Background(), &pb.ListInstanceRequest{ServiceHash: testServiceHash})
		require.NoError(t, err)
		require.Len(t, resp.Instances, 1)
		require.Equal(t, testServiceHash, resp.Instances[0].ServiceHash)
		require.Equal(t, testInstanceHash, resp.Instances[0].Hash)
	})
}

func testDeleteInstance(t *testing.T) {
	ctx := metadata.NewOutgoingContext(context.Background(), passmd)
	_, err := client.InstanceClient.Delete(ctx, &pb.DeleteInstanceRequest{Hash: testInstanceHash})
	require.NoError(t, err)

	resp, err := client.InstanceClient.List(context.Background(), &pb.ListInstanceRequest{ServiceHash: testServiceHash})
	require.NoError(t, err)
	require.Len(t, resp.Instances, 0)
}
