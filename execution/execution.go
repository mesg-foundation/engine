package execution

import (
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/types"
)

// New returns a new execution. It returns an error if inputs are invalid.
func New(processHash, instanceHash, parentHash, eventHash hash.Hash, stepID string, taskKey string, inputs *types.Struct, tags []string) *Execution {
	exec := &Execution{
		ProcessHash:  processHash,
		EventHash:    eventHash,
		InstanceHash: instanceHash,
		ParentHash:   parentHash,
		Inputs:       inputs,
		TaskKey:      taskKey,
		StepId:       stepID,
		Tags:         tags,
		Status:       Status_STATUS_CREATED,
	}
	exec.Hash = hash.Dump(exec)
	return exec
}

// Execute changes executions status to in progres and update its execute time.
// It returns an error if the status is different then Created.
func (execution *Execution) Execute() error {
	if execution.Status != Status_STATUS_CREATED {
		return StatusError{
			ExpectedStatus: Status_STATUS_CREATED,
			ActualStatus:   execution.Status,
		}
	}
	execution.Status = Status_STATUS_IN_PROGRESS
	return nil
}

// Complete changes execution status to completed. It verifies the output.
// It returns an error if the status is different then InProgress or verification fails.
func (execution *Execution) Complete(outputs *types.Struct) error {
	if execution.Status != Status_STATUS_IN_PROGRESS {
		return StatusError{
			ExpectedStatus: Status_STATUS_IN_PROGRESS,
			ActualStatus:   execution.Status,
		}
	}

	execution.Outputs = outputs
	execution.Status = Status_STATUS_COMPLETED
	return nil
}

// Failed changes execution status to failed and puts error information to execution.
// It returns an error if the status is different then InProgress.
func (execution *Execution) Failed(err error) error {
	if execution.Status != Status_STATUS_IN_PROGRESS {
		return StatusError{
			ExpectedStatus: Status_STATUS_IN_PROGRESS,
			ActualStatus:   execution.Status,
		}
	}

	execution.Error = err.Error()
	execution.Status = Status_STATUS_FAILED
	return nil
}
