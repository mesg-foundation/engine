package orchestrator

import (
	"context"

	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/runner"
)

// Store is an interface for fetching all data the orchestrator needs.
type Store interface {
	// FetchProcesses returns all processes.
	FetchProcesses(ctx context.Context) ([]*process.Process, error)

	// FetchExecution returns one execution from its hash.
	FetchExecution(ctx context.Context, hash hash.Hash) (*execution.Execution, error)

	// FetchRunners returns all runners of an instance.
	FetchRunners(ctx context.Context, instanceHash hash.Hash) ([]*runner.Runner, error)

	// CreateExecution creates an execution.
	CreateExecution(ctx context.Context, taskKey string, inputs *types.Struct, tags []string, parentHash hash.Hash, eventHash hash.Hash, processHash hash.Hash, nodeKey string, executorHash hash.Hash) (hash.Hash, error)

	// SubscribeToNewCompletedExecutions returns a chan that will contain newly completed execution.
	SubscribeToNewCompletedExecutions(ctx context.Context) (<-chan *execution.Execution, error)
}
