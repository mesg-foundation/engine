package execution

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/ext/xvalidator"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/tendermint/tendermint/crypto"
)

// New returns a new execution.
func New(processHash, instanceHash, parentHash, eventHash hash.Hash, nodeKey, taskKey, price string, inputs *types.Struct, tags []string, executorHash hash.Hash) (*Execution, error) {
	exec := &Execution{
		ProcessHash:  processHash,
		EventHash:    eventHash,
		InstanceHash: instanceHash,
		ParentHash:   parentHash,
		Inputs:       inputs,
		TaskKey:      taskKey,
		NodeKey:      nodeKey,
		Tags:         tags,
		Price:        price,
		Status:       Status_Proposed,
		ExecutorHash: executorHash,
	}
	exec.Hash = hash.Dump(exec)
	exec.Address = sdk.AccAddress(crypto.AddressHash(exec.Hash))
	return exec, xvalidator.Struct(exec)
}

// Execute changes executions status to in progres.
// It returns an error if the status is different than Proposed.
func (execution *Execution) Execute() error {
	if execution.Status != Status_Proposed {
		return StatusError{
			ExpectedStatus: Status_Proposed,
			ActualStatus:   execution.Status,
		}
	}
	execution.Status = Status_InProgress
	return nil
}

// Complete changes execution status to completed. It verifies the output.
// It returns an error if the status is different than InProgress or verification fails.
func (execution *Execution) Complete(outputs *types.Struct) error {
	if execution.Status != Status_InProgress {
		return StatusError{
			ExpectedStatus: Status_InProgress,
			ActualStatus:   execution.Status,
		}
	}

	execution.Outputs = outputs
	execution.Status = Status_Completed
	return nil
}

// Fail changes execution status to failed and puts error information to execution.
// It returns an error if the status is different than InProgress.
func (execution *Execution) Fail(err error) error {
	if execution.Status != Status_InProgress {
		return StatusError{
			ExpectedStatus: Status_InProgress,
			ActualStatus:   execution.Status,
		}
	}

	execution.Error = err.Error()
	execution.Status = Status_Failed
	return nil
}
