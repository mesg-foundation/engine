package execution

import (
	"github.com/mesg-foundation/engine/hash"
)

// Status stores the state of an execution
type Status int

// Status for an execution
// Created    => The execution is created but not yet processed
// InProgress => The execution is being processed
// Completed  => The execution is completed
const (
	Created Status = iota + 1
	InProgress
	Completed
	Failed
)

func (s Status) String() (r string) {
	switch s {
	case Created:
		r = "created"
	case InProgress:
		r = "in progress"
	case Completed:
		r = "completed"
	case Failed:
		r = "failed"
	}
	return r
}

// Execution stores all information about executions.
type Execution struct {
	Hash         hash.Hash              `hash:"-"`
	WorkflowHash hash.Hash              `hash:"name:workflowHash"`
	ParentHash   hash.Hash              `hash:"name:parentHash"`
	EventHash    hash.Hash              `hash:"name:eventHash"`
	Status       Status                 `hash:"-"`
	InstanceHash hash.Hash              `hash:"name:instanceHash"`
	TaskKey      string                 `hash:"name:taskKey"`
	StepID       string                 `hash:"name:stepID"`
	Tags         []string               `hash:"name:tags"`
	Inputs       map[string]interface{} `hash:"name:inputs"`
	Outputs      map[string]interface{} `hash:"-"`
	Error        string                 `hash:"-"`
}

// New returns a new execution. It returns an error if inputs are invalid.
func New(workflowHash, instanceHash, parentHash, eventHash hash.Hash, stepID string, taskKey string, inputs map[string]interface{}, tags []string) *Execution {
	exec := &Execution{
		WorkflowHash: workflowHash,
		EventHash:    eventHash,
		InstanceHash: instanceHash,
		ParentHash:   parentHash,
		Inputs:       inputs,
		TaskKey:      taskKey,
		StepID:       stepID,
		Tags:         tags,
		Status:       Created,
	}
	exec.Hash = hash.Dump(exec)
	return exec
}

// Execute changes executions status to in progres and update its execute time.
// It returns an error if the status is different then Created.
func (execution *Execution) Execute() error {
	if execution.Status != Created {
		return StatusError{
			ExpectedStatus: Created,
			ActualStatus:   execution.Status,
		}
	}
	execution.Status = InProgress
	return nil
}

// Complete changes execution status to completed. It verifies the output.
// It returns an error if the status is different then InProgress or verification fails.
func (execution *Execution) Complete(outputs map[string]interface{}) error {
	if execution.Status != InProgress {
		return StatusError{
			ExpectedStatus: InProgress,
			ActualStatus:   execution.Status,
		}
	}

	execution.Outputs = outputs
	execution.Status = Completed
	return nil
}

// Failed changes execution status to failed and puts error information to execution.
// It returns an error if the status is different then InProgress.
func (execution *Execution) Failed(err error) error {
	if execution.Status != InProgress {
		return StatusError{
			ExpectedStatus: InProgress,
			ActualStatus:   execution.Status,
		}
	}

	execution.Error = err.Error()
	execution.Status = Failed
	return nil
}
