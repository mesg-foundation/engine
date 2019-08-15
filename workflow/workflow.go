package workflow

import (
	"fmt"

	validator "gopkg.in/go-playground/validator.v9"
)

// Validate returns an error if the workflow is invalid for whatever reason
func (w *Workflow) Validate() error {
	if err := validator.New().Struct(w); err != nil {
		return err
	}
	if w.Trigger.EventKey == "" && w.Trigger.TaskKey == "" {
		return fmt.Errorf("eventKey or taskKey should be present")
	}
	// Check that the initial trigger connects to an existing node.
	if _, err := w.FindNode(w.Trigger.NodeKey); err != nil {
		return err
	}
	// Check that all edges are associated to an existing node.
	for _, edge := range w.Edges {
		if _, err := w.FindNode(edge.Src); err != nil {
			return err
		}
		if _, err := w.FindNode(edge.Dst); err != nil {
			return err
		}
	}
	if err := w.shouldBeDirectedTree(); err != nil {
		return err
	}
	return nil
}

// Match returns true if the data match the current list of filters
func (f TriggerFilters) Match(data map[string]interface{}) bool {
	filters := []*TriggerFilter(f)
	for _, filter := range filters {
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
