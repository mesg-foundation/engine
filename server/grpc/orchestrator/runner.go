package orchestrator

import (
	"context"
	"sync"

	"github.com/mesg-foundation/engine/hash"
)

type runnerStore interface {
	// RegisterRunner registers a new or existing runner.
	RegisterRunner(ctx context.Context, serviceHash hash.Hash, envHash hash.Hash) (hash.Hash, error)

	// DeleteRunner deletes an existing runner.
	DeleteRunner(ctx context.Context, runnerHash hash.Hash) error
}

type runnerServer struct {
	store             runnerStore
	tokenToRunnerHash *sync.Map
	auth              *Authorizer
}

// NewRunnerServer creates a new Runner Server.
func NewRunnerServer(store runnerStore, tokenToRunnerHash *sync.Map, auth *Authorizer) RunnerServer {
	return &runnerServer{
		store:             store,
		tokenToRunnerHash: tokenToRunnerHash,
		auth:              auth,
	}
}

// Register register a new runner.
func (s *runnerServer) Register(ctx context.Context, req *RunnerRegisterRequest) (*RunnerRegisterResponse, error) {
	// check authorization
	if err := s.auth.IsAuthorized(ctx, req); err != nil {
		return nil, err
	}

	// register runner
	runnerHash, err := s.store.RegisterRunner(ctx, req.ServiceHash, req.EnvHash)
	if err != nil {
		return nil, err
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

	return &RunnerRegisterResponse{
		Token: token.String(),
	}, nil
}

// Delete deletes an runner.
func (s *runnerServer) Delete(ctx context.Context, req *RunnerDeleteRequest) (*RunnerDeleteResponse, error) {
	// check authorization
	if err := s.auth.IsAuthorized(ctx, req); err != nil {
		return nil, err
	}

	// delete runner
	if err := s.store.DeleteRunner(ctx, req.RunnerHash); err != nil {
		return nil, err
	}
	return &RunnerDeleteResponse{}, nil
}
