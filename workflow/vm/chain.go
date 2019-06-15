package workflowvm

import (
	"github.com/mesg-foundation/core/workflow"
	uuid "github.com/satori/go.uuid"
)

// executionChain keeps a chain of tasks that should be executed after each other
// while using previous one's outputs as the next one's inputs.
type executionChain struct {
	// w is the workflow of execution chain.
	w workflow.Workflow

	// index of the current task that has been requested to be executed.
	index int
	// indexSecret is a crypto secret that given to current task execution request.
	indexSecret string
}

// newChain creates a new execution chain which is not yet in execution process.
func newChain(w workflow.Workflow) *executionChain {
	return &executionChain{
		w:     w,
		index: -1,
	}
}

// IsDone checks whether is the execution chain can create more execution requests or not.
func (q *executionChain) IsDone() bool {
	return q.index == len(q.w.Tasks)-1
}

// Next creates the next execution requests by using Event as its accumulator.
// Event can be a standard event or a result of previous execution request that in execution chain.
func (q *executionChain) Next(e Event) (ok bool, exec *Execution) {
	if q.IsDone() || // check if there are any more execution requests left in execution chain.
		// check if Event points to the result of the processing execution request of this execution chain.
		// -1 means it is a standard Event, there are no previous execution requests and the first task
		// in the execution chain will be actually put into the processing state.
		q.index != -1 && q.indexSecret != e.Secret {
		return false, nil
	}
	q.index++
	q.indexSecret = uuid.NewV4().String()
	task := q.w.Tasks[q.index]
	return true, &Execution{
		InstanceHash: task.InstanceHash,
		ParentHash:   e.ParentHash,
		TaskKey:      task.Key,
		Inputs:       e.Data,
		Secret:       q.indexSecret,
	}
}
