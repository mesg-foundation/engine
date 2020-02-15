package api

import (
	"context"

	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/protobuf/api"
)

// ProcessServer is the type to aggregate all Service APIs.
type ProcessServer struct {
	mc *cosmos.ModuleClient
}

// NewProcessServer creates a new ProcessServer.
func NewProcessServer(mc *cosmos.ModuleClient) *ProcessServer {
	return &ProcessServer{mc: mc}
}

// Create creates a new process.
func (s *ProcessServer) Create(ctx context.Context, req *api.CreateProcessRequest) (*api.CreateProcessResponse, error) {
	wf, err := s.mc.CreateProcess(req)
	if err != nil {
		return nil, err
	}
	return &api.CreateProcessResponse{Hash: wf.Hash}, nil
}

// Delete deletes process by hash or sid.
func (s *ProcessServer) Delete(ctx context.Context, req *api.DeleteProcessRequest) (*api.DeleteProcessResponse, error) {
	return &api.DeleteProcessResponse{}, s.mc.DeleteProcess(req)
}

// Get returns process from given hash.
func (s *ProcessServer) Get(ctx context.Context, req *api.GetProcessRequest) (*process.Process, error) {
	return s.mc.GetProcess(req.Hash)
}

// List returns all processes.
func (s *ProcessServer) List(ctx context.Context, req *api.ListProcessRequest) (*api.ListProcessResponse, error) {
	processes, err := s.mc.ListProcess()
	if err != nil {
		return nil, err
	}
	return &api.ListProcessResponse{Processes: processes}, nil
}
