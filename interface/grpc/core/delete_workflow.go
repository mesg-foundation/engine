package core

import (
	"context"

	"github.com/mesg-foundation/core/protobuf/coreapi"
)

// DeleteWorkflow stops and deletes a workflow.
func (s *Server) DeleteWorkflow(ctx context.Context, request *coreapi.DeleteWorkflowRequest) (*coreapi.DeleteWorkflowReply, error) {
	return &coreapi.DeleteWorkflowReply{}, s.ss.Workflow().Delete(request.ID)
}
