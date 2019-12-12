package api

import (
	"context"

	protobuf_api "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/runner"
	"github.com/mesg-foundation/engine/sdk"
	runnersdk "github.com/mesg-foundation/engine/sdk/runner"
)

// RunnerServer is the type to aggregate all Runner APIs.
type RunnerServer struct {
	sdk *sdk.SDK
}

// NewRunnerServer creates a new RunnerServer.
func NewRunnerServer(sdk *sdk.SDK) *RunnerServer {
	return &RunnerServer{sdk: sdk}
}

// Create creates a new runner.
func (s *RunnerServer) Create(ctx context.Context, req *protobuf_api.CreateRunnerRequest) (*protobuf_api.CreateRunnerResponse, error) {
	srv, err := s.sdk.Runner.Create(req)
	if err != nil {
		return nil, err
	}
	return &protobuf_api.CreateRunnerResponse{Hash: srv.Hash}, nil
}

// Delete deletes a runner.
func (s *RunnerServer) Delete(ctx context.Context, req *protobuf_api.DeleteRunnerRequest) (*protobuf_api.DeleteRunnerResponse, error) {
	if err := s.sdk.Runner.Delete(req); err != nil {
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
	var filter *runnersdk.Filter
	if req.Filter != nil {
		filter = &runnersdk.Filter{
			Address:      req.Filter.Address,
			InstanceHash: req.Filter.InstanceHash,
		}
	}
	runners, err := s.sdk.Runner.List(filter)
	if err != nil {
		return nil, err
	}

	return &protobuf_api.ListRunnerResponse{Runners: runners}, nil
}
