package main

import (
	"context"
	"testing"

	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/stretchr/testify/require"
)

func testOrchestrator(t *testing.T) {
	var (
		eventStream     pb.Event_StreamClient
		executionStream pb.Execution_StreamClient
		runnerHash      hash.Hash
		instanceHash    hash.Hash
		err             error
	)
	t.Run("setup", func(t *testing.T) {
		respCreate, err := client.ServiceClient.Create(context.Background(), newTestOrchestratorCreateServiceRequest())
		require.NoError(t, err)
		serviceHash := respCreate.Hash

		eventStream, err = client.EventClient.Stream(context.Background(), &pb.StreamEventRequest{})
		require.NoError(t, err)
		acknowledgement.WaitForStreamToBeReady(eventStream)

		executionStream, err = client.ExecutionClient.Stream(context.Background(), &pb.StreamExecutionRequest{})
		require.NoError(t, err)
		acknowledgement.WaitForStreamToBeReady(executionStream)

		respRun, err := client.RunnerClient.Create(context.Background(), &pb.CreateRunnerRequest{
			ServiceHash: serviceHash,
		})
		require.NoError(t, err)
		runnerHash = respRun.Hash

		respRunGet, err := client.RunnerClient.Get(context.Background(), &pb.GetRunnerRequest{Hash: runnerHash})
		require.NoError(t, err)
		instanceHash = respRunGet.InstanceHash

		// wait for service to be ready
		_, err = eventStream.Recv()
		require.NoError(t, err)
	})

	// running orchestrator tests
	t.Run("1 event 1 task", testOrchestrator1Event1Task(executionStream, runnerHash, instanceHash))
	t.Run("1 result 1 task", testOrchestrator1Result1Task(executionStream, runnerHash, instanceHash))

	t.Run("cleanup", func(t *testing.T) {
		_, err = client.RunnerClient.Delete(context.Background(), &pb.DeleteRunnerRequest{Hash: runnerHash})
		require.NoError(t, err)
	})
}
