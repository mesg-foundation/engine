package workflow

import (
	"github.com/mesg-foundation/engine/event"
	"github.com/mesg-foundation/engine/hash"
)

// These structs are temporary and will be part of the service definition
// TODO: Move to service struct

type predicate uint

// Possible status for services.
const (
	EQ predicate = iota
)
type triggerType uint

const (
	// Event is an event emitted by a service
	Event triggerType = iota
	// Result is the result of a task execution
	Result
)

type workflow struct {
	Trigger trigger
	Task    task
}

type task struct {
	InstanceHash hash.Hash
	TaskKey      string
}

type filter struct {
	Key       string
	Predicate predicate
	Value     interface{}
}

// Trigger is an event that triggers a workflow
type trigger struct {
	InstanceHash hash.Hash
	Key          string
	Type         triggerType
	Filters      []*filter
}

func (t *trigger) MatchEvent(evt *event.Event) bool {
	if !t.InstanceHash.Equal(evt.InstanceHash) {
		return false
	}

	if t.Key != evt.Key {
		return false
	}

	for _, filter := range t.Filters {
		if !filter.Match(evt.Data) {
			return false
		}
	}

	return true
}

func (f *filter) Match(inputs map[string]interface{}) bool {
	switch f.Predicate {
	case EQ:
		return inputs[f.Key] == f.Value
	default:
		return false
	}
}
