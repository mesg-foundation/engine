package main

import (
	"context"
	"testing"

	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/stretchr/testify/require"
)

func testOrchestrator(t *testing.T) {
	executionStream, err := client.ExecutionClient.Stream(context.Background(), &pb.StreamExecutionRequest{})
	require.NoError(t, err)
	acknowledgement.WaitForStreamToBeReady(executionStream)

	resultStream, err := client.ResultClient.Stream(context.Background(), &pb.StreamResultRequest{})
	require.NoError(t, err)
	acknowledgement.WaitForStreamToBeReady(resultStream)

	// running orchestrator tests
	t.Run("event task", testOrchestratorEventTask(executionStream, resultStream, testInstanceHash))
	t.Run("result task", testOrchestratorResultTask(executionStream, resultStream, testRunnerHash, testInstanceHash))
	t.Run("event map task", testOrchestratorEventMapTask(executionStream, resultStream, testInstanceHash))
	t.Run("result map task map task", testOrchestratorResultMapTaskMapTask(executionStream, resultStream, testRunnerHash, testInstanceHash))
	t.Run("event map task map task", testOrchestratorEventMapTaskMapTask(executionStream, resultStream, testInstanceHash))
	t.Run("event task complex data", testOrchestratorEventTaskComplexData(executionStream, resultStream, testInstanceHash))
	t.Run("event filter task", testOrchestratorEventFilterTask(executionStream, resultStream, testInstanceHash))
}
