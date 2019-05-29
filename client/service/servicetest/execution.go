package servicetest

import (
	"encoding/json"
)

// Execution is a testing task execution result.
type Execution struct {
	hash string
	data string
	key  string
}

// Hash returns the execution id of task.
func (e *Execution) Hash() string {
	return e.hash
}

// Key returns the output key of task.
func (e *Execution) Key() string {
	return e.key
}

// Data decodes task output to out.
func (e *Execution) Data(out interface{}) error {
	return json.Unmarshal([]byte(e.data), out)
}
