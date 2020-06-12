package orchestrator

import (
	"context"

	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/ext/xstrings"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	types "github.com/mesg-foundation/engine/protobuf/types"
)

// executionStore is the interface to implement to fetch data.
type executionStore interface {
	// CreateExecution creates an execution.
	CreateExecution(ctx context.Context, taskKey string, inputs *types.Struct, tags []string, parentHash hash.Hash, eventHash hash.Hash, processHash hash.Hash, nodeKey string, executorHash hash.Hash) (hash.Hash, error)

	// SubscribeToExecutions returns a chan that will contain executions that have been created, updated, or anything.
	SubscribeToExecutions(ctx context.Context) (<-chan *execution.Execution, error)
}

type executionServer struct {
	store executionStore
	auth  *Authorizer
}

// NewExecutionServer creates a new Execution Server.
func NewExecutionServer(store executionStore, auth *Authorizer) ExecutionServer {
	return &executionServer{
		store: store,
		auth:  auth,
	}
}

// Create creates an execution.
func (s *executionServer) Create(ctx context.Context, req *ExecutionCreateRequest) (*ExecutionCreateResponse, error) {
	// check authorization
	if err := s.auth.IsAuthorized(ctx, req); err != nil {
		return nil, err
	}

	// create execution
	eventHash, err := hash.Random()
	if err != nil {
		return nil, err
	}
	execHash, err := s.store.CreateExecution(
		ctx,
		req.TaskKey,
		req.Inputs,
		req.Tags,
		nil,
		eventHash,
		nil,
		"",
		req.ExecutorHash,
	)
	if err != nil {
		return nil, err
	}
	return &ExecutionCreateResponse{
		Hash: execHash,
	}, nil
}

// Stream returns stream of executions.
func (s *executionServer) Stream(req *ExecutionStreamRequest, stream Execution_StreamServer) error {
	// check authorization
	if err := s.auth.IsAuthorized(stream.Context(), req); err != nil {
		return err
	}

	// create event stream
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()
	executionStream, err := s.store.SubscribeToExecutions(ctx)
	if err != nil {
		return err
	}
	if err := acknowledgement.SetStreamReady(stream); err != nil {
		return err
	}

	// listen to event stream
	for {
		select {
		case exec := <-executionStream:
			if !req.Filter.Match(exec) {
				continue
			}
			if err := stream.Send(exec); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		case <-stream.Context().Done():
			return stream.Context().Err()
		}
	}
}

// Match matches given execution with filter criteria.
func (f *ExecutionStreamRequest_Filter) Match(e *execution.Execution) bool {
	if f == nil {
		return true
	}
	if !f.ExecutorHash.IsZero() && !f.ExecutorHash.Equal(e.ExecutorHash) {
		return false
	}
	if !f.InstanceHash.IsZero() && !f.InstanceHash.Equal(e.InstanceHash) {
		return false
	}
	if f.TaskKey != "" && f.TaskKey != "*" && f.TaskKey != e.TaskKey {
		return false
	}

	match := len(f.Statuses) == 0
	for _, status := range f.Statuses {
		if status == e.Status {
			match = true
			break
		}
	}
	if !match {
		return false
	}
	for _, tag := range f.Tags {
		if !xstrings.SliceContains(e.Tags, tag) {
			return false
		}
	}
	return true
}
