package main

import (
	"context"
	"testing"

	"github.com/mesg-foundation/engine/cosmos/address"
	"github.com/mesg-foundation/engine/ownership"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/service"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
)

var testServiceHash address.testServiceAddress

func testService(t *testing.T) {
	req := newTestCreateServiceRequest()

	t.Run("create", func(t *testing.T) {
		resp, err := client.ServiceClient.Create(context.Background(), req)
		require.NoError(t, err)
		testServiceHash = resp.Hash
	})

	t.Run("get", func(t *testing.T) {
		t.Run("grpc", func(t *testing.T) {
			service, err := client.ServiceClient.Get(context.Background(), &pb.GetServiceRequest{Hash: testServiceHash})
			require.NoError(t, err)
			require.Equal(t, testServiceHash, service.Hash)
		})
		t.Run("lcd", func(t *testing.T) {
			var s *service.Service
			lcdGet(t, "service/get/"+testServiceHash.String(), &s)
			require.Equal(t, testServiceHash, s.Hash)
		})
	})

	t.Run("list", func(t *testing.T) {
		t.Run("grpc", func(t *testing.T) {
			resp, err := client.ServiceClient.List(context.Background(), &pb.ListServiceRequest{})
			require.NoError(t, err)
			require.Len(t, resp.Services, 1)
		})
		t.Run("lcd", func(t *testing.T) {
			ss := make([]*service.Service, 0)
			lcdGet(t, "service/list", &ss)
			require.Len(t, ss, 1)
		})
	})

	t.Run("exists", func(t *testing.T) {
		t.Run("grpc", func(t *testing.T) {
			resp, err := client.ServiceClient.Exists(context.Background(), &pb.ExistsServiceRequest{Hash: testServiceHash})
			require.NoError(t, err)
			require.True(t, resp.Exists)

			resp, err = client.ServiceClient.Exists(context.Background(), &pb.ExistsServiceRequest{Hash: address.ServAddress(crypto.AddressHash([]byte("1")))})
			require.NoError(t, err)
			require.False(t, resp.Exists)
		})
		t.Run("lcd", func(t *testing.T) {
			var exist bool
			lcdGet(t, "service/exist/"+testServiceHash.String(), &exist)
			require.True(t, exist)
			lcdGet(t, "service/exist/"+address.ServAddress(crypto.AddressHash([]byte("1"))).String(), &exist)
			require.False(t, exist)
		})
	})

	t.Run("hash", func(t *testing.T) {
		t.Run("grpc", func(t *testing.T) {
			resp, err := client.ServiceClient.Hash(context.Background(), req)
			require.NoError(t, err)
			require.Equal(t, testServiceHash, resp.Hash)
		})
		t.Run("lcd", func(t *testing.T) {
			var hash address.ServAddress
			lcdPost(t, "service/hash", req, &hash)
			require.Equal(t, testServiceHash, hash)
		})
	})

	t.Run("check ownership creation", func(t *testing.T) {
		t.Run("grpc", func(t *testing.T) {
			ownerships, err := client.OwnershipClient.List(context.Background(), &pb.ListOwnershipRequest{})
			require.NoError(t, err)
			require.Len(t, ownerships.Ownerships, 1)
			require.NotEmpty(t, ownerships.Ownerships[0].Owner)
			require.Equal(t, ownership.Ownership_Service, ownerships.Ownerships[0].Resource)
			require.Equal(t, testServiceHash, ownerships.Ownerships[0].ResourceHash)
		})
		t.Run("lcd", func(t *testing.T) {
			ownerships := make([]*ownership.Ownership, 0)
			lcdGet(t, "ownership/list", &ownerships)
			require.Len(t, ownerships, 1)
			require.NotEmpty(t, ownerships[0].Owner)
			require.Equal(t, ownership.Ownership_Service, ownerships[0].Resource)
			require.Equal(t, testServiceHash, ownerships[0].ResourceHash)
		})
	})
}
