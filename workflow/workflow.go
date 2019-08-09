package workflow

import (
	"fmt"

	"github.com/mesg-foundation/engine/hash"
)

// Match returns true if a workflow trigger is matching the given parameters
func (t *Trigger) Match(trigger TriggerType, instanceHash hash.Hash, key string, data map[string]interface{}) bool {
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
func (f *TriggerFilter) Match(inputs map[string]interface{}) bool {
	switch f.Predicate {
	case EQ:
		return inputs[f.Key] == f.Value
	default:
		return false
	}
}

// ChildrenIDs returns the list of node IDs with a dependency to the current node
func (w Workflow) ChildrenIDs(nodeKey string) []string {
	nodeKeys := make([]string, 0)
	for _, edge := range w.Edges {
		if edge.Src == nodeKey {
			nodeKeys = append(nodeKeys, edge.Dst)
		}
	}
	return nodeKeys
}

// FindNode returns the node matching the key in parameter or an error if not found
func (w Workflow) FindNode(key string) (Node, error) {
	for _, node := range w.Nodes {
		if node.Key == key {
			return node, nil
		}
	}
	return Node{}, fmt.Errorf("%q not found", key)
}
