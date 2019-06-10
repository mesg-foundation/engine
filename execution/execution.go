package execution

import (
	"github.com/mesg-foundation/core/x/xstructhash"
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

// New returns a new execution. It returns an error if inputs are invalid.
func New(service string, parentHash []byte, eventID, taskKey string, inputs map[string]interface{}, tags []string) *Execution {
	exec := &Execution{
		EventID:     eventID,
		ServiceHash: service,
		ParentHash:  parentHash,
		Inputs:      inputs,
		TaskKey:     taskKey,
		Tags:        tags,
		Status:      Created,
	}
	exec.Hash = xstructhash.Hash(exec, 1)
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
