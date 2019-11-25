package main

import (
	"context"
	"testing"

	"github.com/mesg-foundation/engine/hash"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/stretchr/testify/require"
)

var testServiceHash hash.Hash

func testService(t *testing.T) {
	req := newTestCreateServiceRequest()

	t.Run("create", func(t *testing.T) {
		resp, err := client.ServiceClient.Create(context.Background(), req)
		require.NoError(t, err)
		testServiceHash = resp.Hash
	})

	t.Run("get", func(t *testing.T) {
		service, err := client.ServiceClient.Get(context.Background(), &pb.GetServiceRequest{Hash: testServiceHash})
		require.NoError(t, err)
		require.Equal(t, testServiceHash, service.Hash)
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
