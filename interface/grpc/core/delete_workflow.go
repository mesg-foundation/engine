package core

import (
	"context"

	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/systemservices/workflow"
)

// DeleteWorkflow stops and deletes a workflow.
func (s *Server) DeleteWorkflow(ctx context.Context, request *coreapi.DeleteWorkflowRequest) (*coreapi.DeleteWorkflowReply, error) {
	exec, err := s.api.ExecuteAndListen(s.ss.WorkflowServiceID(), workflow.DeleteTaskKey, workflow.DeleteInputs(request.ID))
	if err != nil {
		return nil, err
	}

	return &coreapi.DeleteWorkflowReply{}, workflow.DeleteOutputs(exec)
}
