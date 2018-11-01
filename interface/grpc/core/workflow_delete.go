package core

import (
	"context"

	"github.com/mesg-foundation/core/protobuf/coreapi"
)

// CreateWorkflow creates and runs a new workflow.
func (s *Server) CreateWorkflow(ctx context.Context, request *coreapi.CreateWorkflowRequest) (*coreapi.CreateWorkflowReply, error) {
	id, err := s.ss.Workflow().Create(request.File, request.Name)
	if err != nil {
		return nil, err
	}
	return &coreapi.CreateWorkflowReply{ID: id}, nil
}
