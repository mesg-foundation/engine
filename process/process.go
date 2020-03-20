package process

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/ext/xvalidator"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/tendermint/tendermint/crypto"
)

// New returns a new process and validate it.
func New(name string, nodes []*Process_Node, edges []*Process_Edge) *Process {
	p := &Process{
		Name:  name,
		Nodes: nodes,
		Edges: edges,
	}
	p.Hash = hash.Dump(p)
	p.Address = sdk.AccAddress(crypto.AddressHash(p.Hash))
	return p
}

// Validate returns an error if the process is invalid for whatever reason
func (w *Process) Validate() error {
	if err := xvalidator.Validate.Struct(w); err != nil {
		return err
	}
	if err := w.validate(); err != nil {
		return err
	}
	if _, err := w.Trigger(); err != nil {
		return err
	}
	for _, node := range w.Nodes {
		switch n := node.GetType().(type) {
		case *Process_Node_Map_:
			for _, output := range n.Map.Outputs {
				if ref := output.GetRef(); ref != nil {
					if _, err := w.FindNode(ref.NodeKey); err != nil {
						return err
					}
				}
			}
		case *Process_Node_Filter_:
			for _, condition := range n.Filter.Conditions {
				switch condition.Predicate {
				case Process_Node_Filter_Condition_GT,
					Process_Node_Filter_Condition_GTE,
					Process_Node_Filter_Condition_LT,
					Process_Node_Filter_Condition_LTE:
					if _, ok := condition.Value.Kind.(*types.Value_NumberValue); !ok {
						return fmt.Errorf("filter with condition GT, GTE, LT or LTE only works with value of type Number")
					}
				}
				if _, err := w.FindNode(condition.Ref.NodeKey); err != nil {
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
func (w *Process) Trigger() (*Process_Node, error) {
	triggers := w.FindNodes(func(n *Process_Node) bool {
		return n.GetResult() != nil || n.GetEvent() != nil
	})
	if len(triggers) != 1 {
		return nil, fmt.Errorf("should contain exactly one trigger (result or event)")
	}
	return triggers[0], nil
}
