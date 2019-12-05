package api

import (
	"context"

	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/result"
	"github.com/mesg-foundation/engine/sdk"
)

// ResultServer serve execution functions.
type ResultServer struct {
	sdk *sdk.SDK
}

// NewResultServer creates a new ResultServer.
func NewResultServer(sdk *sdk.SDK) *ResultServer {
	return &ResultServer{sdk: sdk}
}

// Create creates an execution.
func (s *ResultServer) Create(ctx context.Context, req *api.CreateResultRequest) (*api.CreateResultResponse, error) {
	credUsername, credPassphrase, err := GetCredentialFromContext(ctx)
	if err != nil {
		return nil, err
	}
	exec, err := s.sdk.Result.Create(req, credUsername, credPassphrase)
	if err != nil {
		return nil, err
	}

	return &api.CreateResultResponse{
		Hash: exec.Hash,
	}, nil
}

// Get returns execution from given hash.
func (s *ResultServer) Get(ctx context.Context, req *api.GetResultRequest) (*result.Result, error) {
	return s.sdk.Result.Get(req.Hash)
}

// Stream returns stream of executions.
func (s *ResultServer) Stream(req *api.StreamResultRequest, resp api.Result_StreamServer) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, errC, err := s.sdk.Result.Stream(ctx, req)
	if err != nil {
		return err
	}

	if err := acknowledgement.SetStreamReady(resp); err != nil {
		return err
	}

	for {
		select {
		case exec := <-stream:
			if err := resp.Send(exec); err != nil {
				return err
			}
		case err := <-errC:
			return err
		case <-resp.Context().Done():
			return resp.Context().Err()
		}
	}
}

// List returns all executions.
func (s *ResultServer) List(ctx context.Context, req *api.ListResultRequest) (*api.ListResultResponse, error) {
	executions, err := s.sdk.Result.List()
	if err != nil {
		return nil, err
	}
	return &api.ListResultResponse{Results: executions}, nil
}
