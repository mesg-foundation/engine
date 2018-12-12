package service

import (
	"context"

	"github.com/mesg-foundation/core/protobuf/serviceapi"
)

// SubmitResult submits results of an execution.
func (s *Server) SubmitResult(context context.Context, request *serviceapi.SubmitResultRequest) (*serviceapi.SubmitResultReply, error) {
	return &serviceapi.SubmitResultReply{}, s.api.SubmitResult(request.ExecutionID, request.OutputKey, request.OutputData)
}
