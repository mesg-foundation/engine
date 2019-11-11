package api

import (
	"context"
	"errors"

	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/sdk"
	executionsdk "github.com/mesg-foundation/engine/sdk/execution"
)

// ErrNoOutput is an error when there is no output for updating execution.
var ErrNoOutput = errors.New("output not supplied")

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
	eventHash, err := hash.Random()
	if err != nil {
		return nil, err
	}
	executionHash, err := s.sdk.Execution.Execute(nil, req.InstanceHash, eventHash, nil, "", req.TaskKey, req.Inputs, req.Tags)
	if err != nil {
		return nil, err
	}

	return &api.CreateExecutionResponse{
		Hash: executionHash,
	}, nil
}

// Get returns execution from given hash.
func (s *ExecutionServer) Get(ctx context.Context, req *api.GetExecutionRequest) (*execution.Execution, error) {
	return s.sdk.Execution.Get(req.Hash)
}

// Stream returns stream of executions.
func (s *ExecutionServer) Stream(req *api.StreamExecutionRequest, resp api.Execution_StreamServer) error {
	var f *executionsdk.Filter

	if req.Filter != nil {
		f = &executionsdk.Filter{
			InstanceHash: req.Filter.InstanceHash,
			Statuses:     req.Filter.Statuses,
			Tags:         req.Filter.Tags,
			TaskKey:      req.Filter.TaskKey,
		}
	}

	stream := s.sdk.Execution.GetStream(f)
	defer stream.Close()

	// send header to notify client that the stream is ready.
	if err := acknowledgement.SetStreamReady(resp); err != nil {
		return err
	}

	for exec := range stream.C {
		if err := resp.Send(exec); err != nil {
			return err
		}
	}

	return nil
}

// Update updates execution from given hash.
func (s *ExecutionServer) Update(ctx context.Context, req *api.UpdateExecutionRequest) (*api.UpdateExecutionResponse, error) {
	var err error
	switch res := req.Result.(type) {
	case *api.UpdateExecutionRequest_Outputs:
		err = s.sdk.Execution.Update(req.Hash, res.Outputs.Values, nil)
	case *api.UpdateExecutionRequest_Error:
		err = s.sdk.Execution.Update(req.Hash, nil, errors.New(res.Error))
	default:
		err = ErrNoOutput
	}

	if err != nil {
		return nil, err
	}
	return &api.UpdateExecutionResponse{}, nil
}
