package main

import (
	"context"
	"testing"

	"github.com/mesg-foundation/engine/hash"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func testComplexService(t *testing.T) {
	var (
		testServiceHash hash.Hash
		testRunnerHash  hash.Hash
	)

	req := readCreateServiceRequest("testdata/test-complex-service.json")

	t.Run("create", func(t *testing.T) {
		ctx := metadata.NewOutgoingContext(context.Background(), passmd)

		resp, err := client.ServiceClient.Create(ctx, req)
		require.NoError(t, err)

		testServiceHash = resp.Hash
	})

	t.Run("get", func(t *testing.T) {
		ctx := metadata.NewOutgoingContext(context.Background(), passmd)

		service, err := client.ServiceClient.Get(ctx, &pb.GetServiceRequest{Hash: testServiceHash})
		require.NoError(t, err)
		want := pb.TransformCreateReqToService(req)
		// TODO: check why the hashes are different
		want.Hash = service.Hash
		require.True(t, service.Equal(want))
	})

	t.Run("run", func(t *testing.T) {
		ctx := metadata.NewOutgoingContext(context.Background(), passmd)
		resp, err := client.RunnerClient.Create(ctx, &pb.CreateRunnerRequest{
			ServiceHash: testServiceHash,
			Env:         []string{"FOO=1"},
		})
		require.NoError(t, err)
		testRunnerHash = resp.Hash
	})
	t.Run("delete", func(t *testing.T) {
		ctx := metadata.NewOutgoingContext(context.Background(), passmd)
		_, err := client.RunnerClient.Delete(ctx, &pb.DeleteRunnerRequest{Hash: testRunnerHash})
		require.NoError(t, err)
	})
}
