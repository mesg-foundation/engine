package runner

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/event/publisher"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/instance"
	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	"github.com/mesg-foundation/engine/runner"
	executionmodule "github.com/mesg-foundation/engine/x/execution"
	runnermodule "github.com/mesg-foundation/engine/x/runner"
	"google.golang.org/grpc/metadata"
)

// Server is the type to aggregate all Runner APIs.
type Server struct {
	rpc *cosmos.RPC
	ep  *publisher.EventPublisher

	tokenToRunnerHash sync.Map
}

// NewServer creates a new Server.
func NewServer(rpc *cosmos.RPC, ep *publisher.EventPublisher) *Server {
	return &Server{
		rpc: rpc,
		ep:  ep,
	}
}

// Register register a new runner.
func (s *Server) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	// decode msg
	var payload RegisterRequestPayload
	if err := s.rpc.Codec().UnmarshalJSON([]byte(req.Payload), &payload); err != nil {
		return nil, err
	}

	// get account
	acc, err := s.rpc.GetAccount()
	if err != nil {
		return nil, err
	}

	// check signature
	encodedValue, err := s.rpc.Codec().MarshalJSON(payload)
	if err != nil {
		return nil, err
	}
	if !acc.GetPubKey().VerifyBytes(encodedValue, payload.Signature) {
		return nil, fmt.Errorf("verification of the signature failed, it should be signed by %q", acc.GetAddress())
	}

	// calculate runner hash
	inst, err := instance.New(payload.Value.ServiceHash, payload.Value.EnvHash)
	if err != nil {
		return nil, err
	}
	run, err := runner.New(acc.GetAddress().String(), inst.Hash)
	if err != nil {
		return nil, err
	}
	runnerHash := run.Hash

	// check that runner doesn't already exist
	runnerExist := true
	route := fmt.Sprintf("custom/%s/%s/%s", runnermodule.QuerierRoute, runnermodule.QueryGet, runnerHash)
	if _, _, err := s.rpc.QueryWithData(route, nil); err != nil {
		if strings.Contains(err.Error(), fmt.Sprintf("runner %q not found", runnerHash)) {
			runnerExist = false
		} else {
			return nil, err
		}
	}

	// only broadcast if runner doesn't exist
	if !runnerExist {
		tx, err := s.rpc.BuildAndBroadcastMsg(runnermodule.MsgCreate{
			Owner:       acc.GetAddress(),
			ServiceHash: payload.Value.ServiceHash,
			EnvHash:     payload.Value.EnvHash,
		})
		if err != nil {
			return nil, err
		}
		runnerHashCreated, err := hash.DecodeFromBytes(tx.Data)
		if err != nil {
			return nil, err
		}
		if !runnerHashCreated.Equal(runnerHash) {
			// delete wrong runner
			_, err := s.rpc.BuildAndBroadcastMsg(runnermodule.MsgDelete{
				Owner: acc.GetAddress(),
				Hash:  runnerHashCreated,
			})
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("runner hash created is not expected: got %q, expect %q", runnerHashCreated, runnerHash)
		}
	}

	// delete any other token corresponding to runnerHash
	s.tokenToRunnerHash.Range(func(key, value interface{}) bool {
		savedRunnerHash := value.(hash.Hash)
		if savedRunnerHash.Equal(runnerHash) {
			s.tokenToRunnerHash.Delete(key)
		}
		return true
	})

	// generate unique token
	token, err := hash.Random()
	if err != nil {
		return nil, err
	}

	// save token locally with ref to runnerHash
	s.tokenToRunnerHash.Store(token.String(), runnerHash)

	return &RegisterResponse{
		Token: token.String(),
	}, nil
}

// Execution returns a stream of Execution for a specific Runner.
func (s *Server) Execution(req *ExecutionRequest, stream Runner_ExecutionServer) error {
	// check authorization and get runner hash
	runnerHash, err := s.authorize(stream.Context())
	if err != nil {
		return err
	}

	// create rpc event stream
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()
	execChan, execErrChan, err := s.rpc.Stream(ctx, cosmos.EventModuleQuery(executionmodule.ModuleName))
	if err != nil {
		return err
	}
	if err := acknowledgement.SetStreamReady(stream); err != nil {
		return err
	}

	// route channels
	for {
		select {
		case execHash := <-execChan:
			var exec *execution.Execution
			route := fmt.Sprintf("custom/%s/%s/%s", executionmodule.QuerierRoute, executionmodule.QueryGet, execHash)
			if err := s.rpc.QueryJSON(route, nil, &exec); err != nil {
				return err
			}
			// filter execution of this runner
			if exec.ExecutorHash.Equal(runnerHash) && exec.Status == execution.Status_InProgress {
				if err := stream.Send(exec); err != nil {
					return err
				}
			}
		case err := <-execErrChan:
			return err
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
	runnerHash, err := s.authorize(ctx)
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
	msg := executionmodule.MsgUpdate{
		Executor: acc.GetAddress(),
		Hash:     req.ExecutionHash,
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
	return &ResultResponse{}, nil
}

// Event emits an event.
func (s *Server) Event(ctx context.Context, req *EventRequest) (*EventResponse, error) {
	// check authorization and get runner hash
	runnerHash, err := s.authorize(ctx)
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
	if _, err := s.ep.Publish(run.InstanceHash, req.Key, req.Data); err != nil {
		return nil, err
	}

	return &EventResponse{}, nil
}

// --------------------------------------------------------
// Token credential
// --------------------------------------------------------

// TokenCredential is a structure that manage a token.
type TokenCredential struct {
	token string
}

// NewTokenCredential return a token credential struct that implements credentials.PerRPCCredentials interface.
func NewTokenCredential(token string) *TokenCredential {
	return &TokenCredential{
		token: token,
	}
}

// GetRequestMetadata returns the metadata for the request.
func (c *TokenCredential) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"token": c.token,
	}, nil
}

// RequireTransportSecurity tells if the transport should be secured.
func (c *TokenCredential) RequireTransportSecurity() bool {
	return false
	// TODO: test with true
}

// authorize checks the context for a token, matches it against the saved tokens, returns the runner hash if found.
func (s *Server) authorize(ctx context.Context) (hash.Hash, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if len(md["token"]) > 0 {
			token := md["token"][0]
			runnerHash, ok := s.tokenToRunnerHash.Load(token)
			if !ok {
				return nil, fmt.Errorf("credential token doesn't exist")
			}
			return runnerHash.(hash.Hash), nil
		}
	}
	return nil, fmt.Errorf("no credential token found")
}
