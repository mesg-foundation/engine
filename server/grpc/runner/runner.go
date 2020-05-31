package runner

import (
	"context"
	"fmt"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/event/publisher"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/ext/xstrings"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	"github.com/mesg-foundation/engine/runner"
	executionmodule "github.com/mesg-foundation/engine/x/execution"
	runnermodule "github.com/mesg-foundation/engine/x/runner"
	"google.golang.org/grpc/metadata"
)

// CredentialToken is the name to use in the gRPC metadata to set and read the credential token.
const CredentialToken = "mesg_credential_token"

// Server is the type to aggregate all Runner APIs.
type Server struct {
	rpc               *cosmos.RPC
	eventPublisher    *publisher.EventPublisher
	tokenToRunnerHash *sync.Map
	execInProgress    *sync.Map
}

// NewServer creates a new Server.
func NewServer(rpc *cosmos.RPC, eventPublisher *publisher.EventPublisher, tokenToRunnerHash *sync.Map) *Server {
	return &Server{
		rpc:               rpc,
		eventPublisher:    eventPublisher,
		tokenToRunnerHash: tokenToRunnerHash,
		execInProgress:    &sync.Map{},
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

	// create rpc event stream
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()
	subscriber := xstrings.RandASCIILetters(8)
	query := fmt.Sprintf("%s.%s='%s' AND %s.%s='%s'",
		executionmodule.EventType, executionmodule.AttributeKeyExecutor, runnerHash.String(),
		executionmodule.EventType, sdk.AttributeKeyAction, executionmodule.AttributeActionCreated,
	)
	eventStream, err := s.rpc.Subscribe(ctx, subscriber, query, 0)
	defer s.rpc.Unsubscribe(context.Background(), subscriber, query)
	if err != nil {
		return err
	}
	if err := acknowledgement.SetStreamReady(stream); err != nil {
		return err
	}

	// listen to event stream
	for {
		select {
		case event := <-eventStream:
			// get the index of the action=created attributes
			attrKeyActionCreated := fmt.Sprintf("%s.%s", executionmodule.EventType, sdk.AttributeKeyAction)
			attrIndexes := make([]int, 0)
			for index, attr := range event.Events[attrKeyActionCreated] {
				if attr == executionmodule.AttributeActionCreated {
					attrIndexes = append(attrIndexes, index)
				}
			}
			// iterate only on the index of attribute hash where action=created
			attrKeyHash := fmt.Sprintf("%s.%s", executionmodule.EventType, executionmodule.AttributeKeyHash)
			for _, index := range attrIndexes {
				attr := event.Events[attrKeyHash][index]
				hash, err := hash.Decode(attr)
				if err != nil {
					return err
				}
				var exec *execution.Execution
				route := fmt.Sprintf("custom/%s/%s/%s", executionmodule.QuerierRoute, executionmodule.QueryGet, hash)
				if err := s.rpc.QueryJSON(route, nil, &exec); err != nil {
					return err
				}
				s.execInProgress.Store(hash.String(), uint64(time.Now().UnixNano()))
				if err := stream.Send(exec); err != nil {
					return err
				}
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
	var exec *execution.Execution
	route := fmt.Sprintf("custom/%s/%s/%s", executionmodule.QuerierRoute, executionmodule.QueryGet, req.ExecutionHash)
	if err := s.rpc.QueryJSON(route, nil, &exec); err != nil {
		return nil, err
	}
	if !exec.ExecutorHash.Equal(runnerHash) {
		return nil, fmt.Errorf("this runner (%q) is not authorized to submit the result of this execution, the executor should be %q", runnerHash, exec.ExecutorHash)
	}

	// update execution
	acc, err := s.rpc.GetAccount()
	if err != nil {
		return nil, err
	}
	start, ok := s.execInProgress.Load(req.ExecutionHash.String())
	if !ok {
		panic("execution should be in memory")
	}
	msg := executionmodule.MsgUpdate{
		Executor: acc.GetAddress(),
		Hash:     req.ExecutionHash,
		Start:    start.(uint64),
		Stop:     uint64(time.Now().UnixNano()),
	}
	switch result := req.Result.(type) {
	case *ResultRequest_Outputs:
		msg.Result = &executionmodule.MsgUpdateOutputs{
			Outputs: result.Outputs,
		}
	case *ResultRequest_Error:
		msg.Result = &executionmodule.MsgUpdateError{
			Error: result.Error,
		}
	}
	if _, err := s.rpc.BuildAndBroadcastMsg(msg); err != nil {
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
	var run *runner.Runner
	route := fmt.Sprintf("custom/%s/%s/%s", runnermodule.QuerierRoute, runnermodule.QueryGet, runnerHash)
	if err := s.rpc.QueryJSON(route, nil, &run); err != nil {
		return nil, err
	}

	// publish event
	if _, err := s.eventPublisher.Publish(run.InstanceHash, req.Key, req.Data); err != nil {
		return nil, err
	}

	return &EventResponse{}, nil
}
