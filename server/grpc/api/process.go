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
func (s *ProcessServer) Create(ctx context.Context, req *api.ProcessServiceCreateRequest) (*api.ProcessServiceCreateResponse, error) {
	wf, err := s.sdk.Process.Create(&process.Process{
		Key:   req.Key,
		Nodes: req.Nodes,
		Edges: req.Edges,
	})
	if err != nil {
		return nil, err
	}
	return &api.ProcessServiceCreateResponse{Hash: wf.Hash}, nil
}

// Delete deletes service by hash or sid.
func (s *ProcessServer) Delete(ctx context.Context, request *api.ProcessServiceDeleteRequest) (*api.ProcessServiceDeleteResponse, error) {
	return &api.ProcessServiceDeleteResponse{}, s.sdk.Process.Delete(request.Hash)
}

// Get returns service from given hash.
func (s *ProcessServer) Get(ctx context.Context, req *api.ProcessServiceGetRequest) (*api.ProcessServiceGetResponse, error) {
	process, err := s.sdk.Process.Get(req.Hash)
	if err != nil {
		return nil, err
	}
	return &api.ProcessServiceGetResponse{Process: process}, nil
}

// List returns all processes.
func (s *ProcessServer) List(ctx context.Context, req *api.ProcessServiceListRequest) (*api.ProcessServiceListResponse, error) {
	processes, err := s.sdk.Process.List()
	if err != nil {
		return nil, err
	}
	return &api.ProcessServiceListResponse{
		Processes: processes,
	}, nil
}
