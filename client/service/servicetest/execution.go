package servicetest

import "encoding/json"

// Execution is a testing task execution result.
type Execution struct {
	id   string
	data string
}

// ID returns the execution id of task.
func (e *Execution) ID() string {
	return e.id
}

// Data decodes task output to out.
func (e *Execution) Data(out interface{}) error {
	return json.Unmarshal([]byte(e.data), out)
}
