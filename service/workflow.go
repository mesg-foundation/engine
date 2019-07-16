package service

import "github.com/mesg-foundation/engine/hash"

// Match returns true if a workflow trigger is matching the given parameters
func (t *WorkflowTrigger) Match(trigger TriggerType, instanceHash hash.Hash, key string, data map[string]interface{}) bool {
	if t.Type != trigger {
		return false
	}
	if !t.InstanceHash.Equal(instanceHash) {
		return false
	}

	if t.Key != key {
		return false
	}

	for _, filter := range t.Filters {
		if !filter.Match(data) {
			return false
		}
	}

	return true
}

// Match returns true the current filter matches the given data
func (f *WorkflowTriggerFilter) Match(inputs map[string]interface{}) bool {
	switch f.Predicate {
	case EQ:
		return inputs[f.Key] == f.Value
	default:
		return false
	}
}
