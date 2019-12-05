package api

import (
	fmt "fmt"

	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/x/xstrings"
)

// Validate checks if given filter is valid and returns error.
func (f *StreamExecutionRequest_Filter) Validate() error {
	if f == nil {
		return nil
	}

	if !f.ExecutorHash.Valid() {
		return fmt.Errorf("stream filter: executor hash is invalid")
	}

	if !f.InstanceHash.Valid() {
		return fmt.Errorf("stream filter: instance hash is invalid")
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
	if !f.ExecutorHash.IsZero() && !f.ExecutorHash.Equal(e.ExecutorHash) {
		return false
	}
	if !f.InstanceHash.IsZero() && !f.InstanceHash.Equal(e.InstanceHash) {
		return false
	}
	if f.TaskKey != "" && f.TaskKey != "*" && f.TaskKey != e.TaskKey {
		return false
	}
	for _, tag := range f.Tags {
		if !xstrings.SliceContains(e.Tags, tag) {
			return false
		}
	}
	return true
}
