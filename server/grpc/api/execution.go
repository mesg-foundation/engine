package api

import (
	"context"

	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/sdk"
)

// ExecutionServer serve execution functions.
type ExecutionServer struct {
	sdk *sdk.SDK
}

// NewExecutionServer creates a new ExecutionServer.
func NewExecutionServer(sdk *sdk.SDK) *ExecutionServer {
	return &ExecutionServer{sdk: sdk}
}

// Create creates an execution.
func (s *ExecutionServer) Create(ctx context.Context, req *api.CreateExecutionRequest) (*api.CreateExecutionResponse, error) {
	credUsername, credPassphrase, err := GetCredentialFromContext(ctx)
	if err != nil {
		return nil, err
	}
	exec, err := s.sdk.Execution.Create(req, credUsername, credPassphrase)
	if err != nil {
		return nil, err
	}

	return &api.CreateExecutionResponse{
		Hash: exec.Hash,
	}, nil
}

// Get returns execution from given hash.
func (s *ExecutionServer) Get(ctx context.Context, req *api.GetExecutionRequest) (*execution.Execution, error) {
	return s.sdk.Execution.Get(req.Hash)
}

// Stream returns stream of executions.
func (s *ExecutionServer) Stream(req *api.StreamExecutionRequest, resp api.Execution_StreamServer) error {
	stream, err := s.sdk.Execution.Stream(req)
	if err != nil {
		return err
	}
	defer close(stream)

	if err := acknowledgement.SetStreamReady(resp); err != nil {
		return err
	}

	// TODO:
	// There is possible deadlock. If the client close the connection,
	// but there will be no messages in the stream, then this for will
	// wait and consume resources forever. Some ACK mechnizm needs to be
	// implemented on server/client side to get notify if the conneciton
	// wasn't closed.
	for exec := range stream {
		if err := resp.Send(exec); err != nil {
			return err
		}
	}
	return nil
}

// Update updates execution from given hash.
func (s *ExecutionServer) Update(ctx context.Context, req *api.UpdateExecutionRequest) (*api.UpdateExecutionResponse, error) {
	credUsername, credPassphrase, err := GetCredentialFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := s.sdk.Execution.Update(req, credUsername, credPassphrase); err != nil {
		return nil, err
	}
	return &api.UpdateExecutionResponse{}, nil
}

// List returns all executions.
func (s *ExecutionServer) List(ctx context.Context, req *api.ListExecutionRequest) (*api.ListExecutionResponse, error) {
	executions, err := s.sdk.Execution.List()
	if err != nil {
		return nil, err
	}
	return &api.ListExecutionResponse{Executions: executions}, nil
}
