package main

import (
	"testing"
)

func testOrchestrator(t *testing.T) {
	// running orchestrator tests
	t.Run("process balance and withdraw", testOrchestratorProcessBalanceWithdraw(testRunnerHash, testInstanceHash))
	t.Run("event task", testOrchestratorEventTask(testRunnerHash, testInstanceHash))
	t.Run("result task", testOrchestratorResultTask(testRunnerHash, testInstanceHash))
	t.Run("map const", testOrchestratorMapConst(testRunnerHash, testInstanceHash))
	t.Run("ref grand parent task", testOrchestratorRefGrandParentTask(testRunnerHash, testInstanceHash))
	t.Run("nested data", testOrchestratorNestedData(testRunnerHash, testInstanceHash))
	t.Run("nested map", testOrchestratorNestedMap(testRunnerHash, testInstanceHash))
	t.Run("ref path nested", testOrchestratorRefPathNested(testRunnerHash, testInstanceHash))
	t.Run("filter", testOrchestratorFilter(testRunnerHash, testInstanceHash))
	t.Run("filter path nested", testOrchestratorFilterPathNested(testRunnerHash, testInstanceHash))
}
