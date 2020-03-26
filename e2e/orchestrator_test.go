package main

import (
	"testing"
)

func testOrchestrator(t *testing.T) {
	// running orchestrator tests
	t.Run("process balance and withdraw", testOrchestratorProcessBalanceWithdraw(testInstanceHash))
	t.Run("event task", testOrchestratorEventTask(testInstanceHash))
	t.Run("result task", testOrchestratorResultTask(testRunnerHash, testInstanceHash))
	t.Run("map const", testOrchestratorMapConst(testInstanceHash))
	t.Run("ref grand parent task", testOrchestratorRefGrandParentTask(testInstanceHash))
	t.Run("nested data", testOrchestratorNestedData(testInstanceHash))
	t.Run("nested map", testOrchestratorNestedMap(testInstanceHash))
	t.Run("ref path nested", testOrchestratorRefPathNested(testInstanceHash))

	// to execute last because of go routine leak. See fixme in following function
	t.Run("filter", testOrchestratorFilter(testInstanceHash))
	t.Run("filter path nested", testOrchestratorFilterPathNested(testInstanceHash))
}
