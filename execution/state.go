package execution

import (
	"errors"
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

func (execution *Execution) moveFromPendingToInProgress() (err error) {
	mu.Lock()
	defer mu.Unlock()
	e, ok := pendingExecutions[execution.ID]
	if !ok {
		err = errors.New("Execution " + execution.ID + " not in the pending")
		return
	}
	inProgressExecutions[execution.ID] = e
	delete(pendingExecutions, execution.ID)
	return
}

func (execution *Execution) moveFromInProgressToProcessed() (err error) {
	mu.Lock()
	defer mu.Unlock()
	e, ok := inProgressExecutions[execution.ID]
	if !ok {
		err = errors.New("Execution " + execution.ID + " not in progress")
		return
	}
	processedExecutions[execution.ID] = e
	delete(inProgressExecutions, execution.ID)
	return
}
