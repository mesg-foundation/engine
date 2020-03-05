package api

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	execution "github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/ext/xstrings"
)

// Validate checks if given filter is valid and returns error.
func (f *StreamExecutionRequest_Filter) Validate() error {
	if f == nil {
		return nil
	}

	if !f.ExecutorHash.Empty() {
		if err := sdk.VerifyAddressFormat(f.ExecutorHash); err != nil {
			return fmt.Errorf("stream filter: executor hash is invalid: %w", err)
		}
	}

	if !f.InstanceHash.Empty() {
		if err := sdk.VerifyAddressFormat(f.InstanceHash); err != nil {
			return fmt.Errorf("stream filter: instance hash is invalid: %w", err)
		}
	}

	// TODO: add validation (after adding in protobuf with print ascii)
	// if f.TaskKey == "" || f.TaskKey == "*" || validation {
	// return err
	// }

	return nil
}

// Match matches given execution with filter criteria.
func (f *StreamExecutionRequest_Filter) Match(e *execution.Execution) bool {
	if f == nil {
		return true
	}
	if !f.ExecutorHash.Empty() && !f.ExecutorHash.Equals(e.ExecutorHash) {
		return false
	}
	if !f.InstanceHash.Empty() && !f.InstanceHash.Equals(e.InstanceHash) {
		return false
	}
	if f.TaskKey != "" && f.TaskKey != "*" && f.TaskKey != e.TaskKey {
		return false
	}

	match := len(f.Statuses) == 0
	for _, status := range f.Statuses {
		if status == e.Status {
			match = true
			break
		}
	}
	if !match {
		return false
	}
	for _, tag := range f.Tags {
		if !xstrings.SliceContains(e.Tags, tag) {
			return false
		}
	}
	return true
}
