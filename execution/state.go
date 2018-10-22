package execution

import (
	"sync"
)

var mu sync.Mutex
var pendingExecutions = make(map[string]*Execution)
var inProgressExecutions = make(map[string]*Execution)

// InProgress returns the matching in progress execution if exists
func InProgress(ID string) (execution *Execution) {
	mu.Lock()
	defer mu.Unlock()
	execution = inProgressExecutions[ID]
	return execution
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

func (execution *Execution) deleteFromInProgressQueue() error {
	mu.Lock()
	defer mu.Unlock()
	_, ok := inProgressExecutions[execution.ID]
	if !ok {
		return &NotInQueueError{
			ID:    execution.ID,
			Queue: "inProgress",
		}
	}
	delete(inProgressExecutions, execution.ID)
	return nil
}
