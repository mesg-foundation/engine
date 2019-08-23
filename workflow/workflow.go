package workflow

import (
	"fmt"

	validator "gopkg.in/go-playground/validator.v9"
)

// ID is the ID of the Result's node
func (r Result) ID() string {
	return "trigger:result"
}

// ID is the ID of the Event's node
func (r Event) ID() string {
	return "trigger:event"
}

// ID is the ID of the Task's node
func (m Task) ID() string {
	return m.Key
}

// ID is the ID of the Mapping's node
func (m Mapping) ID() string {
	return m.Key
}

// Validate returns an error if the workflow is invalid for whatever reason
func (w *Workflow) Validate() error {
	if err := validator.New().Struct(w); err != nil {
		return err
	}
	if err := w.Graph.validate(); err != nil {
		return err
	}
	triggers := w.Graph.FindNodes(func(n Node) bool {
		_, isResult := n.(Result)
		_, isEvent := n.(Event)
		return isResult || isEvent
	})
	if len(triggers) != 1 {
		return fmt.Errorf("should contain exactly one trigger (result or event)")
	}
	for _, node := range w.Graph.Nodes {
		switch n := node.(type) {
		case Mapping:
			for _, input := range n.Inputs {
				if _, err := w.FindNode(input.Ref.NodeKey); err != nil {
					return err
				}
			}
		}
	}
	if err := w.shouldBeDirectedTree(); err != nil {
		return err
	}
	return nil
}
