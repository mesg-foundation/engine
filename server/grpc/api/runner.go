package api

import (
	"context"

	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/runner"
	"github.com/mesg-foundation/engine/runner/builder"
)

// RunnerServer is the type to aggregate all Runner APIs.
type RunnerServer struct {
	mc *cosmos.ModuleClient
	b  *builder.Builder
}

// NewRunnerServer creates a new RunnerServer.
func NewRunnerServer(mc *cosmos.ModuleClient, b *builder.Builder) *RunnerServer {
	return &RunnerServer{
		mc: mc,
		b:  b,
	}
}

// Create creates a new runner.
func (s *RunnerServer) Create(ctx context.Context, req *api.CreateRunnerRequest) (*api.CreateRunnerResponse, error) {
	r, err := s.b.Create(req)
	if err != nil {
		return nil, err
	}
	return &api.CreateRunnerResponse{Hash: r.Hash}, nil
}

// Delete deletes a runner.
func (s *RunnerServer) Delete(ctx context.Context, req *api.DeleteRunnerRequest) (*api.DeleteRunnerResponse, error) {
	if err := s.b.Delete(req); err != nil {
		return nil, err
	}
	return &api.DeleteRunnerResponse{}, nil
}

// Get returns runner from given hash.
func (s *RunnerServer) Get(ctx context.Context, req *api.GetRunnerRequest) (*runner.Runner, error) {
	return s.mc.GetRunner(req.Hash)
}

// List returns all runners.
func (s *RunnerServer) List(ctx context.Context, req *api.ListRunnerRequest) (*api.ListRunnerResponse, error) {
	var f *cosmos.FilterRunner
	if req.Filter != nil {
		f = &cosmos.FilterRunner{
			Address:      req.Filter.Address,
			InstanceHash: req.Filter.InstanceHash,
		}
	}
	runners, err := s.mc.ListRunner(f)
	if err != nil {
		return nil, err
	}

	return &api.ListRunnerResponse{Runners: runners}, nil
}
