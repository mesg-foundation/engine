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
	return execution.move("pending", pendingExecutions, inProgressExecutions)
}

func (execution *Execution) moveFromInProgressToProcessed() error {
	return execution.move("inProgress", inProgressExecutions, processedExecutions)
}

func (execution *Execution) move(queue string, from, to map[string]*Execution) error {
	mu.Lock()
	defer mu.Unlock()
	e, ok := from[execution.ID]
	if !ok {
		return &NotInQueueError{
			ID:    execution.ID,
			Queue: queue,
		}
	}
	to[execution.ID] = e
	delete(from, execution.ID)
	return nil
}
