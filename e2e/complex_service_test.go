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

func testComplexService(t *testing.T) {
	var (
		testServiceHash  hash.Hash
		testRunnerHash   hash.Hash
		testInstanceHash hash.Hash
	)

	req := newTestComplexCreateServiceRequest()
	ctx := metadata.NewOutgoingContext(context.Background(), passmd)

	t.Run("create", func(t *testing.T) {
		resp, err := client.ServiceClient.Create(ctx, req)
		require.NoError(t, err)

		testServiceHash = resp.Hash
	})

	stream, err := client.EventClient.Stream(context.Background(), &pb.StreamEventRequest{})
	require.NoError(t, err)
	acknowledgement.WaitForStreamToBeReady(stream)

	t.Run("run", func(t *testing.T) {
		resp, err := client.RunnerClient.Create(ctx, &pb.CreateRunnerRequest{
			ServiceHash: testServiceHash,
			Env:         []string{"FOO=bar"},
		})
		require.NoError(t, err)
		testRunnerHash = resp.Hash

		resp1, err := client.RunnerClient.Get(context.Background(), &pb.GetRunnerRequest{Hash: testRunnerHash})
		require.NoError(t, err)
		testInstanceHash = resp1.InstanceHash
	})
	t.Run("check events", func(t *testing.T) {
		okEventsNo := 4
		for i := 0; i < okEventsNo; {
			ev, err := stream.Recv()
			require.NoError(t, err)

			if !ev.InstanceHash.Equal(testInstanceHash) {
				continue
			}
			i++

			switch ev.Key {
			case "test_service_ready", "read_env_ok",
				"access_volumes_from_ok", "ping_dependencie_ok":
				t.Logf("recieved event %s ", ev.Key)
			default:
				t.Fatalf("failed on event %s", ev.Key)
			}
		}
	})

	t.Run("delete", func(t *testing.T) {
		_, err := client.RunnerClient.Delete(ctx, &pb.DeleteRunnerRequest{Hash: testRunnerHash})
		require.NoError(t, err)
	})
}
