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
	t.Run("map const", testOrchestratorMapConst(executionStream, testInstanceHash))
	t.Run("ref grand parent task", testOrchestratorRefGrandParentTask(executionStream, testInstanceHash))
	t.Run("nested data", testOrchestratorNestedData(executionStream, testInstanceHash))
	t.Run("nested map", testOrchestratorNestedMap(executionStream, testInstanceHash))
	t.Run("ref path nested", testOrchestratorRefPathNested(executionStream, testInstanceHash))

	// to execute last because of go routine leak. See fixme in following function
	t.Run("filter", testOrchestratorFilter(executionStream, testInstanceHash))
}
