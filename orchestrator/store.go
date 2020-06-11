package orchestrator

import (
	"context"

	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/runner"
)

type Store interface {
	FetchProcesses(ctx context.Context) ([]*process.Process, error)
	FetchExecution(ctx context.Context, hash hash.Hash) (*execution.Execution, error)
	FetchRunners(ctx context.Context, instanceHash hash.Hash) ([]*runner.Runner, error)
	CreateExecution(ctx context.Context, taskKey string, inputs *types.Struct, tags []string, parentHash hash.Hash, eventHash hash.Hash, processHash hash.Hash, nodeKey string, executorHash hash.Hash) (hash.Hash, error)
	SubscribeToNewCompletedExecutions(ctx context.Context) (<-chan *execution.Execution, error)
}
