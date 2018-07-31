package mesgtest

import "encoding/json"

// Execution is a testing task execution result.
type Execution struct {
	id   string
	data string
	key  string
}

// ID returns the execution id of task.
func (e *Execution) ID() string {
	return e.id
}

// Key returns the output key of task.
func (e *Execution) Key() string {
	return e.key
}

// Decode decodes task output to out.
func (e *Execution) Decode(out interface{}) error {
	return json.Unmarshal([]byte(e.data), out)
}
