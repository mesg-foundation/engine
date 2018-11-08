package core

import (
	"context"

	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/systemservices/workflow"
)

// CreateWorkflow creates and runs a new workflow.
func (s *Server) CreateWorkflow(ctx context.Context, request *coreapi.CreateWorkflowRequest) (*coreapi.CreateWorkflowReply, error) {
	exec, err := s.api.ExecuteAndListen(s.ss.WorkflowServiceID(), workflow.CreateTaskKey, workflow.CreateInputs(request.File, request.Name))
	if err != nil {
		return nil, err
	}

	id, err := workflow.CreateOutputs(exec)
	if err != nil {
		return nil, err
	}
	return &coreapi.CreateWorkflowReply{ID: id}, nil
}
