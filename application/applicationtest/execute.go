package applicationtest

import "encoding/json"

// Execute holds information about a task execution.
type Execute struct {
	serviceID string
	task      string
	data      string
}

// ServiceID returns the id of service that task executed on.
func (e *Execute) ServiceID() string {
	return e.serviceID
}

// Task returns the executed task's name.
func (e *Execute) Task() string {
	return e.task
}

// Decode decodes task's input data to out.
func (e *Execute) Decode(out interface{}) error {
	return json.Unmarshal([]byte(e.data), out)
}
