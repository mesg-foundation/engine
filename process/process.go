package process

import (
	"fmt"

	validator "gopkg.in/go-playground/validator.v9"
)

// ID is the ID of the Result's node
func (r Result) ID() string {
	return r.Key
}

// ID is the ID of the Event's node
func (e Event) ID() string {
	return e.Key
}

// ID is the ID of the Task's node
func (m Task) ID() string {
	return m.Key
}

// ID is the ID of the Map's node
func (m Map) ID() string {
	return m.Key
}

// ID is the ID of the Filter's node
func (f Filter) ID() string {
	return f.Key
}

// Validate returns an error if the process is invalid for whatever reason
func (w *Process) Validate() error {
	if err := validator.New().Struct(w); err != nil {
		return err
	}
	if err := w.Graph.validate(); err != nil {
		return err
	}
	if _, err := w.Trigger(); err != nil {
		return err
	}
	for _, node := range w.Graph.Nodes {
		n, isMap := node.(Map)
		if isMap {
			for _, output := range n.Outputs {
				if _, err := w.FindNode(output.Ref.NodeKey); err != nil {
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

// Trigger returns the trigger of the process
func (w *Process) Trigger() (Node, error) {
	triggers := w.Graph.FindNodes(func(n Node) bool {
		_, isResult := n.(Result)
		_, isEvent := n.(Event)
		return isResult || isEvent
	})
	if len(triggers) != 1 {
		return nil, fmt.Errorf("should contain exactly one trigger (result or event)")
	}
	return triggers[0], nil
}