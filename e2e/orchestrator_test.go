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

	// running orchestrator tests
	t.Run("event task", testOrchestratorEventTask(executionStream, testInstanceHash))
	t.Run("result task", testOrchestratorResultTask(executionStream, testRunnerHash, testInstanceHash))
	t.Run("event map task", testOrchestratorEventMapTask(executionStream, testInstanceHash))
	t.Run("result map task map task", testOrchestratorResultMapTaskMapTask(executionStream, testRunnerHash, testInstanceHash))
	t.Run("event task map task map task", testOrchestratorEventTaskMapTaskMapTask(executionStream, testInstanceHash))
	t.Run("event task complex data", testOrchestratorEventTaskComplexData(executionStream, testInstanceHash))
	t.Run("event map task complex", testOrchestratorEventMapTaskComplex(executionStream, testInstanceHash))
	t.Run("event map task map task complex data", testOrchestratorEventMapTaskMapTaskComplexData(executionStream, testInstanceHash))

	// to execute last because of go routine leak. See fixme in following function
	t.Run("event filter task", testOrchestratorEventFilterTask(executionStream, testInstanceHash))
}
