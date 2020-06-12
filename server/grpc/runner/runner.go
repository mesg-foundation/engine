package runner

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/mesg-foundation/engine/event/publisher"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	types "github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/runner"
	tmlog "github.com/tendermint/tendermint/libs/log"
	"google.golang.org/grpc/metadata"
)

// Store is the interface to implement to fetch data.
type Store interface {
	// SubscribeToExecutionsForRunner returns a chan that will contain executions that a specific runner must execute.
	SubscribeToExecutionsForRunner(ctx context.Context, runnerHash hash.Hash) (<-chan *execution.Execution, error)

	// FetchExecution returns one execution from its hash.
	FetchExecution(ctx context.Context, hash hash.Hash) (*execution.Execution, error)

	// UpdateExecution update an execution.
	UpdateExecution(ctx context.Context, execHash hash.Hash, start int64, stop int64, outputs *types.Struct, err string) error

	// FetchRunner returns a runner from its hash.
	FetchRunner(ctx context.Context, hash hash.Hash) (*runner.Runner, error)
}

// CredentialToken is the name to use in the gRPC metadata to set and read the credential token.
const CredentialToken = "mesg_credential_token"

// Server is the type to aggregate all Runner APIs.
type Server struct {
	store             Store
	eventPublisher    *publisher.EventPublisher
	tokenToRunnerHash *sync.Map
	execInProgress    *sync.Map
	logger            tmlog.Logger
}

// NewServer creates a new Server.
func NewServer(store Store, eventPublisher *publisher.EventPublisher, tokenToRunnerHash *sync.Map, logger tmlog.Logger) *Server {
	return &Server{
		store:             store,
		eventPublisher:    eventPublisher,
		tokenToRunnerHash: tokenToRunnerHash,
		execInProgress:    &sync.Map{},
		logger:            logger,
	}
}

// isAuthorized checks the context for a token, matches it against the saved tokens, returns the runner hash if found.
func (s *Server) isAuthorized(ctx context.Context) (hash.Hash, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("credential not found in metadata, make sure to set it using the name %q", CredentialToken)
	}
	if len(md[CredentialToken]) == 0 {
		return nil, fmt.Errorf("credential not found in metadata, make sure to set it using the name %q", CredentialToken)
	}
	token := md[CredentialToken][0]
	runnerHash, ok := s.tokenToRunnerHash.Load(token)
	if !ok {
		return nil, fmt.Errorf("credential token doesn't exist")
	}
	return runnerHash.(hash.Hash), nil
}

// Execution returns a stream of Execution for a specific Runner.
func (s *Server) Execution(req *ExecutionRequest, stream Runner_ExecutionServer) error {
	// check authorization and get runner hash
	runnerHash, err := s.isAuthorized(stream.Context())
	if err != nil {
		return err
	}

	// create event stream
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()
	executionStream, err := s.store.SubscribeToExecutionsForRunner(ctx, runnerHash)
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
			s.execInProgress.Store(exec.Hash.String(), time.Now().UnixNano())
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

// Result emits the result of an Execution.
func (s *Server) Result(ctx context.Context, req *ResultRequest) (*ResultResponse, error) {
	// check authorization and get runner hash
	runnerHash, err := s.isAuthorized(ctx)
	if err != nil {
		return nil, err
	}

	// make sure runner is allowed to update this execution
	exec, err := s.store.FetchExecution(ctx, req.ExecutionHash)
	if err != nil {
		return nil, err
	}
	if !exec.ExecutorHash.Equal(runnerHash) {
		return nil, fmt.Errorf("this runner (%q) is not authorized to submit the result of this execution, the executor should be %q", runnerHash, exec.ExecutorHash)
	}

	// update execution
	start, ok := s.execInProgress.Load(req.ExecutionHash.String())
	if !ok {
		s.logger.Error(fmt.Sprintf("execution %q should be in memory", req.ExecutionHash.String()))
		start = time.Now().UnixNano()
	}
	if err := s.store.UpdateExecution(
		ctx,
		req.ExecutionHash,
		start.(int64),
		time.Now().UnixNano(),
		req.GetOutputs(),
		req.GetError(),
	); err != nil {
		return nil, err
	}
	s.execInProgress.Delete(req.ExecutionHash.String())
	return &ResultResponse{}, nil
}

// Event emits an event.
func (s *Server) Event(ctx context.Context, req *EventRequest) (*EventResponse, error) {
	// check authorization and get runner hash
	runnerHash, err := s.isAuthorized(ctx)
	if err != nil {
		return nil, err
	}

	// get runner to access instance hash
	run, err := s.store.FetchRunner(ctx, runnerHash)

	// publish event
	if _, err := s.eventPublisher.Publish(run.InstanceHash, req.Key, req.Data); err != nil {
		return nil, err
	}

	return &EventResponse{}, nil
}
