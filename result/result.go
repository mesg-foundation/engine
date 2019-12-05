package result

import (
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/types"
)

// NewWithOutputs returns a new result
func NewWithOutputs(requestHash hash.Hash, outputs *types.Struct) *Result {
	result := &Result{
		RequestHash: requestHash,
		Result:      &Result_Outputs{outputs},
	}
	result.Hash = hash.Dump(result)
	return result
}

// NewWithError returns a new result
func NewWithError(requestHash hash.Hash, err string) *Result {
	result := &Result{
		RequestHash: requestHash,
		Result:      &Result_Error{err},
	}
	result.Hash = hash.Dump(result)
	return result
}
