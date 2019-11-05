package api

import (
	"context"

	"github.com/mesg-foundation/engine/config"
	protobuf_api "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/runner"
	"github.com/mesg-foundation/engine/sdk"
)

// RunnerServer is the type to aggregate all Runner APIs.
type RunnerServer struct {
	sdk *sdk.SDK
	cfg *config.Config
}

// NewRunnerServer creates a new RunnerServer.
func NewRunnerServer(sdk *sdk.SDK, cfg *config.Config) *RunnerServer {
	return &RunnerServer{sdk: sdk, cfg: cfg}
}

// Create creates a new runner.
func (s *RunnerServer) Create(ctx context.Context, req *protobuf_api.CreateRunnerRequest) (*protobuf_api.CreateRunnerResponse, error) {
	// credUsername, credPassphrase, err := GetCredentialFromContext(ctx)
	// if err != nil {
	// 	return nil, err
	// }
	credUsername, credPassphrase := s.cfg.Account.Name, s.cfg.Account.Password

	srv, err := s.sdk.Runner.Create(req, credUsername, credPassphrase)
	if err != nil {
		return nil, err
	}
	return &protobuf_api.CreateRunnerResponse{Hash: srv.Hash}, nil
}

// Delete deletes a runner.
func (s *RunnerServer) Delete(ctx context.Context, req *protobuf_api.DeleteRunnerRequest) (*protobuf_api.DeleteRunnerResponse, error) {
	// credUsername, credPassphrase, err := GetCredentialFromContext(ctx)
	// if err != nil {
	// 	return nil, err
	// }
	credUsername, credPassphrase := s.cfg.Account.Name, s.cfg.Account.Password

	if err := s.sdk.Runner.Delete(req, credUsername, credPassphrase); err != nil {
		return nil, err
	}
	return &protobuf_api.DeleteRunnerResponse{}, nil
}

// Get returns runner from given hash.
func (s *RunnerServer) Get(ctx context.Context, req *protobuf_api.GetRunnerRequest) (*runner.Runner, error) {
	return s.sdk.Runner.Get(req.Hash)
}

// List returns all runners.
func (s *RunnerServer) List(ctx context.Context, req *protobuf_api.ListRunnerRequest) (*protobuf_api.ListRunnerResponse, error) {
	runners, err := s.sdk.Runner.List(req.Filter)
	if err != nil {
		return nil, err
	}

	return &protobuf_api.ListRunnerResponse{Runners: runners}, nil
}
