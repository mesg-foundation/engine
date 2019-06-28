package api

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/hash"
	"github.com/mesg-foundation/core/protobuf/acknowledgement"
	"github.com/mesg-foundation/core/protobuf/api"
	"github.com/mesg-foundation/core/protobuf/types"
	"github.com/mesg-foundation/core/sdk"
	executionsdk "github.com/mesg-foundation/core/sdk/execution"
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

// Get returns execution from given hash.
func (s *ExecutionServer) Get(ctx context.Context, req *api.GetExecutionRequest) (*types.Execution, error) {
	hash, err := hash.Decode(req.Hash)
	if err != nil {
		return nil, err
	}

	exec, err := s.sdk.Execution.Get(hash)
	if err != nil {
		return nil, err
	}
	return toProtoExecution(exec)
}

// Stream returns stream of executions.
func (s *ExecutionServer) Stream(req *api.StreamExecutionRequest, resp api.Execution_StreamServer) error {
	var f *executionsdk.Filter

	if req.Filter != nil {
		instanceHash, err := hash.Decode(req.Filter.InstanceHash)
		if req.Filter.InstanceHash != "" && err != nil {
			return err
		}

		var statuses []execution.Status
		for _, status := range req.Filter.Statuses {
			statuses = append(statuses, execution.Status(status))
		}

		f = &executionsdk.Filter{
			InstanceHash: instanceHash,
			Statuses:     statuses,
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
func (s *ExecutionServer) Update(ctx context.Context, req *api.UpdateExecutionRequest) (*api.UpdateExecutionResponse, error) {
	hash, err := hash.Decode(req.Hash)
	if err != nil {
		return nil, err
	}
	switch res := req.Result.(type) {
	case *api.UpdateExecutionRequest_Outputs:
		err = s.sdk.Execution.Update(hash, []byte(res.Outputs), nil)
	case *api.UpdateExecutionRequest_Error:
		err = s.sdk.Execution.Update(hash, nil, errors.New(res.Error))
	default:
		err = ErrNoOutput
	}

	if err != nil {
		return nil, err
	}
	return &api.UpdateExecutionResponse{}, nil

}

func toProtoExecution(exec *execution.Execution) (*types.Execution, error) {
	inputs, err := json.Marshal(exec.Inputs)
	if err != nil {
		return nil, err
	}

	outputs, err := json.Marshal(exec.Outputs)
	if err != nil {
		return nil, err
	}

	return &types.Execution{
		Hash:         exec.Hash.String(),
		ParentHash:   exec.ParentHash.String(),
		EventID:      exec.EventID,
		Status:       types.Status(exec.Status),
		InstanceHash: exec.InstanceHash.String(),
		TaskKey:      exec.TaskKey,
		Inputs:       string(inputs),
		Outputs:      string(outputs),
		Tags:         exec.Tags,
		Error:        exec.Error,
	}, nil
}
