package api

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/protobuf/acknowledgement"
	"github.com/mesg-foundation/core/protobuf/api"
	"github.com/mesg-foundation/core/protobuf/definition"
	"github.com/mesg-foundation/core/sdk"
)

// ErrNoOutput is an error when there is no output for updating execution.
var ErrNoOutput = errors.New("output not supplied")

// Server serve execution functions.
type Server struct {
	sdk *sdk.SDK
}

// NewServer creates a new Server.
func NewServer(sdk *sdk.SDK) *Server {
	return &Server{sdk: sdk}
}

// Get returns execution from given hash.
func (s *Server) Get(ctx context.Context, req *api.GetExecutionRequest) (*definition.Execution, error) {
	hash, err := hex.DecodeString(req.Hash)
	if err != nil {
		return nil, err
	}

	exec, err := s.sdk.GetExecution(hash)
	if err != nil {
		return nil, err
	}
	return toProtoExecution(exec)
}

// Stream returns stream of executions.
func (s *Server) Stream(req *api.StreamExecutionRequest, resp api.Execution_StreamServer) error {
	stream := s.sdk.GetExecutionStream(&sdk.ExecutionFilter{
		Statuses: []execution.Status{execution.Status(req.Filter.Status)},
	})
	defer stream.Close()

	// send header to notify client that the stream is ready.
	if err := acknowledgement.SetStreamReady(resp); err != nil {
		return err
	}

	for exec := range stream.C {
		pexec, err := toProtoExecution(exec)
		if err != nil {
			return err
		}

		if err := resp.Send(pexec); err != nil {
			return err
		}
	}

	return nil
}

// Update updates execution from given hash.
func (s *Server) Update(ctx context.Context, req *api.UpdateExecutionRequest) (*api.UpdateExecutionResponse, error) {
	hash, err := hex.DecodeString(req.Hash)
	if err != nil {
		return nil, err
	}
	switch res := req.Result.(type) {
	case *api.UpdateExecutionRequest_Outputs:
		err = s.sdk.UpdateExecution(hash, []byte(res.Outputs), nil)
	case *api.UpdateExecutionRequest_Error:
		err = s.sdk.UpdateExecution(hash, nil, errors.New(res.Error))
	default:
		err = ErrNoOutput
	}

	if err != nil {
		return nil, err
	}
	return &api.UpdateExecutionResponse{}, nil

}

func toProtoExecution(exec *execution.Execution) (*definition.Execution, error) {
	inputs, err := json.Marshal(exec.Inputs)
	if err != nil {
		return nil, err
	}

	outputs, err := json.Marshal(exec.Outputs)
	if err != nil {
		return nil, err
	}

	return &definition.Execution{
		Hash:        hex.EncodeToString(exec.Hash),
		ParentHash:  hex.EncodeToString(exec.ParentHash),
		EventID:     exec.EventID,
		Status:      definition.Status(exec.Status),
		ServiceHash: exec.ServiceHash,
		TaskKey:     exec.TaskKey,
		Inputs:      string(inputs),
		Outputs:     string(outputs),
		Tags:        exec.Tags,
		Error:       exec.Error,
	}, nil
}
