package runnersdk

import (
	"fmt"

	"github.com/mesg-foundation/engine/hash"
)

// AlreadyExistsError is an not found error.
type AlreadyExistsError struct {
	Hash hash.Hash
}

func (e *AlreadyExistsError) Error() string {
	return fmt.Sprintf("runner %q already exists", e.Hash.String())
}
