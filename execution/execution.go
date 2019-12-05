package execution

import (
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/types"
)

// New returns a new execution
func New(processHash, instanceHash, parentResultHash, eventHash hash.Hash, nodeKey string, taskKey string, inputs *types.Struct, tags []string, executorHash hash.Hash) *Execution {
	exec := &Execution{
		ProcessHash:      processHash,
		EventHash:        eventHash,
		InstanceHash:     instanceHash,
		ParentResultHash: parentResultHash,
		Inputs:           inputs,
		TaskKey:          taskKey,
		NodeKey:          nodeKey,
		Tags:             tags,
		ExecutorHash:     executorHash,
	}
	exec.Hash = hash.Dump(exec)
	return exec
}
