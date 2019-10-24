package main

import (
	"context"
	"testing"

	"github.com/mesg-foundation/engine/hash"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func testInstance(t *testing.T) {
	var instanceHash hash.Hash
	t.Run("create", func(t *testing.T) {
		ctx := metadata.NewOutgoingContext(context.Background(), passmd)
		resp, err := client.InstanceClient.Create(ctx, &pb.CreateInstanceRequest{ServiceHash: testServiceHash})
		require.NoError(t, err)
		instanceHash = resp.Hash
	})

	t.Run("get", func(t *testing.T) {
		resp, err := client.InstanceClient.Get(context.Background(), &pb.GetInstanceRequest{Hash: instanceHash})
		require.NoError(t, err)
		require.Equal(t, instanceHash, resp.Hash)
		require.Equal(t, testServiceHash, resp.ServiceHash)
	})

	t.Run("list", func(t *testing.T) {
		resp, err := client.InstanceClient.List(context.Background(), &pb.ListInstanceRequest{ServiceHash: testServiceHash})
		require.NoError(t, err)
		require.Len(t, resp.Instances, 1)
	})

	// t.Run("delete", func(t *testing.T) {
	// 	ctx := metadata.NewOutgoingContext(context.Background(), passmd)
	// 	_, err := client.InstanceClient.Delete(ctx, &pb.DeleteInstanceRequest{Hash: instanceHash})
	// 	require.NoError(t, err)
	// })
}
