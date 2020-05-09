package orchestrator

func (l *OrchestratorLog) keyvals() []interface{} {
	return []interface{}{
		"processHash", l.ProcessHash.String(),
		"nodeKey", l.NodeKey,
		"nodeType", l.NodeType,
		"timeConsumed", l.TimeConsumed,
	}
}

type orchestratorLogType interface {
	isOrchestratorLog_Type
	message() string
	keyvals() []interface{}
}

func (a *OrchestratorLog_ExecutionCreated) message() string {
	return "execution created"
}
func (a *OrchestratorLog_ExecutionCreated) keyvals() []interface{} {
	return []interface{}{
		"executionHash", a.ExecutionCreated.ExecutionHash.String(),
	}
}

func (a *OrchestratorLog_MapExecuted) message() string {
	return "map executed"
}
func (a *OrchestratorLog_MapExecuted) keyvals() []interface{} {
	return []interface{}{}
}

func (a *OrchestratorLog_FilterMatched) message() string {
	return "filter matched"
}
func (a *OrchestratorLog_FilterMatched) keyvals() []interface{} {
	return []interface{}{}
}

func (a *OrchestratorLog_FilterDidNotMatched) message() string {
	return "filter did not matched"
}
func (a *OrchestratorLog_FilterDidNotMatched) keyvals() []interface{} {
	return []interface{}{}
}

func (a *OrchestratorLog_TriggeredByResult) message() string {
	return "triggered by result"
}
func (a *OrchestratorLog_TriggeredByResult) keyvals() []interface{} {
	return []interface{}{
		"executionHash", a.TriggeredByResult.ExecutionHash.String(),
	}
}

func (a *OrchestratorLog_TriggeredByEvent) message() string {
	return "triggered by event"
}
func (a *OrchestratorLog_TriggeredByEvent) keyvals() []interface{} {
	return []interface{}{
		"eventHash", a.TriggeredByEvent.EventHash.String(),
	}
}

func (a *OrchestratorLog_ErrorOccurred) message() string {
	return "an error occurred"
}
func (a *OrchestratorLog_ErrorOccurred) keyvals() []interface{} {
	return []interface{}{
		"error", a.ErrorOccurred.Error,
	}
}
