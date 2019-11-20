package main

import (
	"context"
	"testing"

	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/stretchr/testify/require"
)

func testComplexService(t *testing.T) {
	var (
		testServiceHash  hash.Hash
		testRunnerHashC  hash.Hash
		testInstanceHash hash.Hash
	)

	req := newTestComplexCreateServiceRequest()

	t.Run("create", func(t *testing.T) {
		resp, err := client.ServiceClient.Create(context.Background(), req)
		require.NoError(t, err)
		testServiceHash = resp.Hash
	})

	stream, err := client.EventClient.Stream(context.Background(), &pb.StreamEventRequest{})
	require.NoError(t, err)
	acknowledgement.WaitForStreamToBeReady(stream)

	t.Run("run", func(t *testing.T) {
		resp, err := client.RunnerClient.Create(context.Background(), &pb.CreateRunnerRequest{
			ServiceHash: testServiceHash,
			Env:         []string{"ENVB=is_override"},
		})
		require.NoError(t, err)
		testRunnerHashC = resp.Hash

		resp1, err := client.RunnerClient.Get(context.Background(), &pb.GetRunnerRequest{Hash: testRunnerHashC})
		require.NoError(t, err)
		testInstanceHash = resp1.InstanceHash
	})

	t.Run("check events", func(t *testing.T) {
		okEventsNo := 6
		for i := 0; i < okEventsNo; {
			ev, err := stream.Recv()
			require.NoError(t, err)

			if !ev.InstanceHash.Equal(testInstanceHash) {
				continue
			}
			i++

			switch ev.Key {
			case "test_service_ready", "read_env_ok", "read_env_override_ok", "access_volumes_ok", "access_volumes_from_ok", "resolve_dependence_ok":
				t.Logf("received event %s ", ev.Key)
			default:
				t.Fatalf("failed on event %s", ev.Key)
			}
		}
	})

	t.Run("delete", func(t *testing.T) {
		t.Skip("FIXME: this call never get trough. some issue with the service's dependency")
		_, err := client.RunnerClient.Delete(context.Background(), &pb.DeleteRunnerRequest{Hash: testRunnerHashC})
		require.NoError(t, err)
	})
}
