package process

import (
	"fmt"

	validator "gopkg.in/go-playground/validator.v9"
)

// ID is the ID of the Result's node
func (r *Process_Node) ID() string {
	if event := r.GetEvent(); event != nil {
		return event.Key
	}
	if filter := r.GetFilter(); filter != nil {
		return filter.Key
	}
	if mapping := r.GetMap(); mapping != nil {
		return mapping.Key
	}
	if result := r.GetResult(); result != nil {
		return result.Key
	}
	if task := r.GetTask(); task != nil {
		return task.Key
	}
	panic("not implemented")
}

// Validate returns an error if the process is invalid for whatever reason
func (w *Process) Validate() error {
	if err := validator.New().Struct(w); err != nil {
		return err
	}
	if err := w.validate(); err != nil {
		return err
	}
	if _, err := w.Trigger(); err != nil {
		return err
	}
	for _, node := range w.Nodes {
		mapNode := node.GetMap()
		if mapNode != nil {
			for _, output := range mapNode.Outputs {
				if ref := output.GetRef(); ref != nil {
					if _, err := w.FindNode(ref.NodeKey); err != nil {
						return err
					}
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
func (w *Process) Trigger() (*Process_Node, error) {
	triggers := w.FindNodes(func(n *Process_Node) bool {
		return n.GetResult() != nil || n.GetEvent() != nil
	})
	if len(triggers) != 1 {
		return nil, fmt.Errorf("should contain exactly one trigger (result or event)")
	}
	return triggers[0], nil
}
