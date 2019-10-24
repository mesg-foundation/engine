package main

import (
	"context"
	"testing"

	"github.com/mesg-foundation/engine/hash"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func testService(t *testing.T) {
	const serviceHash = "4F1Nst7TeKVbWXa3So9WZRmQEcqwwF8PEVkBLXyCAAqa"

	req := readCreateServiceRequest("testdata/test-service/compile.json")

	t.Run("create", func(t *testing.T) {
		ctx := metadata.NewOutgoingContext(context.Background(), passmd)

		resp, err := client.ServiceClient.Create(ctx, req)
		require.NoError(t, err)
		require.Equal(t, serviceHash, resp.Hash.String())
		testServiceHash = resp.Hash
	})

	t.Run("get", func(t *testing.T) {
		ctx := metadata.NewOutgoingContext(context.Background(), passmd)

		service, err := client.ServiceClient.Get(ctx, &pb.GetServiceRequest{Hash: testServiceHash})
		require.NoError(t, err)
		require.Equal(t, serviceHash, service.Hash.String())
	})

	t.Run("list", func(t *testing.T) {
		resp, err := client.ServiceClient.List(context.Background(), &pb.ListServiceRequest{})
		require.NoError(t, err)
		require.Len(t, resp.Services, 1)
	})

	t.Run("exists", func(t *testing.T) {
		resp, err := client.ServiceClient.Exists(context.Background(), &pb.ExistsServiceRequest{Hash: testServiceHash})
		require.NoError(t, err)
		require.True(t, resp.Exists)

		resp, err = client.ServiceClient.Exists(context.Background(), &pb.ExistsServiceRequest{Hash: hash.Int(1)})
		require.NoError(t, err)
		require.False(t, resp.Exists)
	})

	t.Run("hash", func(t *testing.T) {
		resp, err := client.ServiceClient.Hash(context.Background(), req)
		require.NoError(t, err)
		require.Equal(t, testServiceHash, resp.Hash)
	})
}
