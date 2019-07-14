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
	EventKey     string
	Filters      []*filter
}

func (t *trigger) Match(evt *event.Event) (bool, error) {
	instanceHash := evt.InstanceHash
	if !t.InstanceHash.Equal(instanceHash) {
		return false, nil
	}

	if t.EventKey != evt.Key {
		return false, nil
	}

	for _, filter := range t.Filters {
		if !filter.Match(evt.Data) {
			return false, nil
		}
	}

	return true, nil
}

func (f *filter) Match(inputs map[string]interface{}) bool {
	switch f.Predicate {
	case EQ:
		return inputs[f.Key] == f.Value
	default:
		return false
	}
}
