package orchestrator

import (
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/event"
	"github.com/mesg-foundation/engine/event/publisher"
	"github.com/mesg-foundation/engine/execution"
)

// Orchestrator manages the executions based on the definition of the processes
type Orchestrator struct {
	rpc *cosmos.RPC
	ep  *publisher.EventPublisher

	eventStream *event.Listener

	executionStream chan *execution.Execution
	ErrC            chan error
	stopC           chan bool

	accName     string
	accPassword string
	execPrice   string
}
