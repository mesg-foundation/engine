package api

import (
	"context"

	"github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/sdk"
)

// ProcessServer is the type to aggregate all Service APIs.
type ProcessServer struct {
	sdk *sdk.SDK
}

// NewProcessServer creates a new ProcessServer.
func NewProcessServer(sdk *sdk.SDK) *ProcessServer {
	return &ProcessServer{sdk: sdk}
}

// Create creates a new service from definition.
func (s *ProcessServer) Create(ctx context.Context, req *api.CreateProcessRequest) (*api.CreateProcessResponse, error) {
	wf, err := s.sdk.Process.Create(&process.Process{
		Key:   req.Key,
		Nodes: req.Nodes,
		Edges: req.Edges,
	})
	if err != nil {
		return nil, err
	}
	return &api.CreateProcessResponse{Hash: wf.Hash}, nil
}

// Delete deletes service by hash or sid.
func (s *ProcessServer) Delete(ctx context.Context, request *api.DeleteProcessRequest) (*api.DeleteProcessResponse, error) {
	return &api.DeleteProcessResponse{}, s.sdk.Process.Delete(request.Hash)
}

// Get returns service from given hash.
func (s *ProcessServer) Get(ctx context.Context, req *api.GetProcessRequest) (*process.Process, error) {
	return s.sdk.Process.Get(req.Hash)
}

// List returns all processes.
func (s *ProcessServer) List(ctx context.Context, req *api.ListProcessRequest) (*api.ListProcessResponse, error) {
	processes, err := s.sdk.Process.List()
	if err != nil {
		return nil, err
	}
	return &api.ListProcessResponse{
		Processes: processes,
	}, nil
}
