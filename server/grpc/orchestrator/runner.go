package orchestrator

import (
	"context"
	fmt "fmt"
	"sync"

	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/instance"
	"github.com/mesg-foundation/engine/runner"
	runnermodule "github.com/mesg-foundation/engine/x/runner"
)

type runnerServer struct {
	rpc               *cosmos.RPC
	tokenToRunnerHash *sync.Map
	auth              *Authorizer
}

// NewRunnerServer creates a new Runner Server.
func NewRunnerServer(rpc *cosmos.RPC, tokenToRunnerHash *sync.Map, auth *Authorizer) RunnerServer {
	return &runnerServer{
		rpc:               rpc,
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

	// get engine account
	acc, err := s.rpc.GetAccount()
	if err != nil {
		return nil, err
	}

	// calculate runner hash
	inst, err := instance.New(req.ServiceHash, req.EnvHash)
	if err != nil {
		return nil, err
	}
	run, err := runner.New(acc.GetAddress().String(), inst.Hash)
	if err != nil {
		return nil, err
	}
	runnerHash := run.Hash

	// check that runner doesn't already exist
	var runnerExist bool
	route := fmt.Sprintf("custom/%s/%s/%s", runnermodule.QuerierRoute, runnermodule.QueryExist, runnerHash)
	if err := s.rpc.QueryJSON(route, nil, &runnerExist); err != nil {
		return nil, err
	}

	// only broadcast if runner doesn't exist
	if !runnerExist {
		tx, err := s.rpc.BuildAndBroadcastMsg(runnermodule.MsgCreate{
			Owner:       acc.GetAddress(),
			ServiceHash: req.ServiceHash,
			EnvHash:     req.EnvHash,
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

	// create execution
	acc, err := s.rpc.GetAccount()
	if err != nil {
		return nil, err
	}
	msg := runnermodule.MsgDelete{
		Owner: acc.GetAddress(),
		Hash:  req.RunnerHash,
	}
	if _, err := s.rpc.BuildAndBroadcastMsg(msg); err != nil {
		return nil, err
	}
	return &RunnerDeleteResponse{}, nil
}
