package orchestrator

import (
	"context"

	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/runner"
	eventsdk "github.com/mesg-foundation/engine/sdk/event"
	runnersdk "github.com/mesg-foundation/engine/sdk/runner"
)

// ExecutionSDK execution interface needed for the orchestrator
type ExecutionSDK interface {
	Stream(ctx context.Context, req *api.StreamExecutionRequest) (chan *execution.Execution, chan error, error)
	Get(hash hash.Hash) (*execution.Execution, error)
	Create(req *api.CreateExecutionRequest) (*execution.Execution, error)
}

// EventSDK event interface needed for the orchestrator
type EventSDK interface {
	GetStream(f *eventsdk.Filter) *eventsdk.Listener
}

// ProcessSDK process interface needed for the orchestrator
type ProcessSDK interface {
	List() ([]*process.Process, error)
}

// RunnerSDK is the interface of the runner sdk needed for the orchestrator
type RunnerSDK interface {
	List(f *runnersdk.Filter) ([]*runner.Runner, error)
}

// Orchestrator manages the executions based on the definition of the processes
type Orchestrator struct {
	event       EventSDK
	eventStream *eventsdk.Listener

	execution       ExecutionSDK
	executionStream <-chan *execution.Execution

	process ProcessSDK

	runner RunnerSDK

	ErrC chan error
}
