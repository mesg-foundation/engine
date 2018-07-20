package execution

import (
	"sync"
)

var mu sync.Mutex
var pendingExecutions = make(map[string]*Execution)
var inProgressExecutions = make(map[string]*Execution)
var processedExecutions = make(map[string]*Execution)

// InProgress returns the matching in progress execution if exists
func InProgress(ID string) (execution *Execution) {
	execution = inProgressExecutions[ID]
	return
}

func (execution *Execution) moveFromPendingToInProgress() error {
	mu.Lock()
	defer mu.Unlock()
	e, ok := pendingExecutions[execution.ID]
	if !ok {
		return &NotInQueueError{
			ID:    execution.ID,
			Queue: "pending",
		}
	}
	inProgressExecutions[execution.ID] = e
	delete(pendingExecutions, execution.ID)
	return nil
}

func (execution *Execution) moveFromInProgressToProcessed() error {
	mu.Lock()
	defer mu.Unlock()
	e, ok := inProgressExecutions[execution.ID]
	if !ok {
		return &NotInQueueError{
			ID:    execution.ID,
			Queue: "inProgress",
		}
	}
	processedExecutions[execution.ID] = e
	delete(inProgressExecutions, execution.ID)
	return nil
}
