package process

import (
	"fmt"

	validator "gopkg.in/go-playground/validator.v9"
)

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
				if ref := output.GetRef(); ref != nil && ref.RefKey != "" {
					refNode, err := w.FindTaskNodeFromRef(node.Key, ref.RefKey)
					if err != nil {
						return err
					}
					if refNode == nil {
						return fmt.Errorf("%q reference not found", ref.RefKey)
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

// FindTaskNodeFromRef returns the node associated to a reference in the nodeKey's parents, returns an nil if not found
func (w *Process) FindTaskNodeFromRef(nodeKey string, refKey string) (*Process_Node, error) {
	node, err := w.FindNode(nodeKey)
	if err != nil {
		return nil, err
	}
	if node.GetTask() != nil && node.GetTask().RefKey == refKey {
		return node, nil
	}
	var refNode *Process_Node
	for _, parent := range w.ParentKeys(nodeKey) {
		r, err := w.FindTaskNodeFromRef(parent, refKey)
		if err != nil {
			return nil, err
		}
		if refNode == nil && r != nil {
			refNode = r
		}
	}
	return refNode, nil
}
