package result

import (
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/types"
)

// NewWithOutputs returns a new result
func NewWithOutputs(executionHash hash.Hash, outputs *types.Struct) *Result {
	result := &Result{
		ExecutionHash: executionHash,
		Result:        &Result_Outputs{outputs},
	}
	result.Hash = hash.Dump(result)
	return result
}

// NewWithError returns a new result
func NewWithError(executionHash hash.Hash, err string) *Result {
	result := &Result{
		ExecutionHash: executionHash,
		Result:        &Result_Error{err},
	}
	result.Hash = hash.Dump(result)
	return result
}
