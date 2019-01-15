package mesgtest

import "encoding/json"

// Event represents a testing event.
type Event struct {
	name  string
	data  string
	token string
}

// Name returns the name of event.
func (e *Event) Name() string {
	return e.name
}

// Data decodes event data to out.
func (e *Event) Data(out interface{}) error {
	return json.Unmarshal([]byte(e.data), out)
}

// Token returns the service id of service emitted the event.
func (e *Event) Token() string {
	return e.token
}
