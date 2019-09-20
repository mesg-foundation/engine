package orchestrator

import (
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/protobuf/types"
	eventsdk "github.com/mesg-foundation/engine/sdk/event"
	executionsdk "github.com/mesg-foundation/engine/sdk/execution"
)

// ExecutionSDK execution interface needed for the orchestrator
type ExecutionSDK interface {
	GetStream(f *executionsdk.Filter) *executionsdk.Listener
	Get(hash hash.Hash) (*execution.Execution, error)
	Execute(processHash, instanceHash, eventHash, parentHash hash.Hash, stepID string, taskKey string, inputData *types.Struct, tags []string) (executionHash hash.Hash, err error)
}

// EventSDK event interface needed for the orchestrator
type EventSDK interface {
	GetStream(f *eventsdk.Filter) *eventsdk.Listener
}

// ProcessSDK process interface needed for the orchestrator
type ProcessSDK interface {
	List() ([]*process.Process, error)
}

// Orchestrator manages the executions based on the definition of the processes
type Orchestrator struct {
	event       EventSDK
	eventStream *eventsdk.Listener

	execution       ExecutionSDK
	executionStream *executionsdk.Listener

	process ProcessSDK

	ErrC chan error
}
