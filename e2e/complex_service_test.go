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

	req := newTestComplexCreateServiceRequest()
	ctx := metadata.NewOutgoingContext(context.Background(), passmd)

	t.Run("create", func(t *testing.T) {
		resp, err := client.ServiceClient.Create(ctx, req)
		require.NoError(t, err)

		testServiceHash = resp.Hash
	})

	t.Run("run", func(t *testing.T) {
		resp, err := client.RunnerClient.Create(ctx, &pb.CreateRunnerRequest{
			ServiceHash: testServiceHash,
			Env:         []string{"FOO=1"},
		})
		require.NoError(t, err)
		testRunnerHash = resp.Hash
	})

	t.Run("delete", func(t *testing.T) {
		_, err := client.RunnerClient.Delete(ctx, &pb.DeleteRunnerRequest{Hash: testRunnerHash})
		require.NoError(t, err)
	})
}
