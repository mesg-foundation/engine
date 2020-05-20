package types

import (
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
)

// Querier names
const (
	QueryGet  = "get"
	QueryList = "list"
)

// ListFilter available for the List
type ListFilter struct {
	ParentHash   hash.Hash        `json:"parentHash"`
	EventHash    hash.Hash        `json:"eventHash"`
	InstanceHash hash.Hash        `json:"instanceHash"`
	ProcessHash  hash.Hash        `json:"processHash"`
	Status       execution.Status `json:"status"`
	NodeKey      string           `json:"nodeKey"`
}

// Match returns true if an execution matches a specific filter
func (f ListFilter) Match(exec *execution.Execution) bool {
	if f.Status != execution.Status_Unknown && f.Status != exec.Status {
		return false
	}
	if !f.ProcessHash.IsZero() && !f.ProcessHash.Equal(exec.ProcessHash) {
		return false
	}
	if !f.InstanceHash.IsZero() && !f.InstanceHash.Equal(exec.InstanceHash) {
		return false
	}
	if !f.ParentHash.IsZero() && !f.ParentHash.Equal(exec.ParentHash) {
		return false
	}
	if !f.EventHash.IsZero() && !f.EventHash.Equal(exec.EventHash) {
		return false
	}
	if f.NodeKey != "" && f.NodeKey != exec.NodeKey {
		return false
	}
	return true
}
