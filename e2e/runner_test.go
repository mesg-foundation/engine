package main

import (
	"context"
	"testing"

	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/runner"
	"github.com/stretchr/testify/require"
)

var testRunnerHash address.testRunnerAddress

func testRunner(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		stream, err := client.EventClient.Stream(context.Background(), &pb.StreamEventRequest{
			Filter: &pb.StreamEventRequest_Filter{
				Key: "test_service_ready",
			},
		})
		require.NoError(t, err)
		acknowledgement.WaitForStreamToBeReady(stream)

		resp, err := client.RunnerClient.Create(context.Background(), &pb.CreateRunnerRequest{
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
		t.Run("grpc", func(t *testing.T) {
			resp, err := client.RunnerClient.Get(context.Background(), &pb.GetRunnerRequest{Hash: testRunnerHash})
			require.NoError(t, err)
			require.Equal(t, testRunnerHash, resp.Hash)
			testInstanceHash = resp.InstanceHash
		})
		t.Run("lcd", func(t *testing.T) {
			var r *runner.Runner
			lcdGet(t, "runner/get/"+testRunnerHash.String(), &r)
			require.Equal(t, testRunnerHash, r.Hash)
		})
	})

	// TODO: need to test the filters
	t.Run("list", func(t *testing.T) {
		t.Run("grpc", func(t *testing.T) {
			resp, err := client.RunnerClient.List(context.Background(), &pb.ListRunnerRequest{})
			require.NoError(t, err)
			require.Len(t, resp.Runners, 1)
			require.Equal(t, testInstanceHash, resp.Runners[0].InstanceHash)
			require.Equal(t, testRunnerHash, resp.Runners[0].Hash)
		})
		t.Run("lcd", func(t *testing.T) {
			rs := make([]*runner.Runner, 0)
			lcdGet(t, "runner/list", &rs)
			require.Len(t, rs, 1)
			require.Equal(t, testInstanceHash, rs[0].InstanceHash)
			require.Equal(t, testRunnerHash, rs[0].Hash)
		})
	})
}

func testDeleteRunner(t *testing.T) {
	_, err := client.RunnerClient.Delete(context.Background(), &pb.DeleteRunnerRequest{Hash: testRunnerHash})
	require.NoError(t, err)

	t.Run("grpc", func(t *testing.T) {
		resp, err := client.RunnerClient.List(context.Background(), &pb.ListRunnerRequest{})
		require.NoError(t, err)
		require.Len(t, resp.Runners, 0)
	})
	t.Run("lcd", func(t *testing.T) {
		rs := make([]*runner.Runner, 0)
		lcdGet(t, "runner/list", &rs)
		require.Len(t, rs, 0)
	})
}
