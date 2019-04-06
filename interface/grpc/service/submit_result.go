package service

import (
	"context"
	"fmt"

	"github.com/mesg-foundation/core/protobuf/serviceapi"
)

// SubmitResult submits results of an execution.
func (s *Server) SubmitResult(context context.Context, request *serviceapi.SubmitResultRequest) (*serviceapi.SubmitResultReply, error) {
	var err error
	if request.GetError() != "" {
		err = fmt.Errorf(request.GetError())
	}
	return &serviceapi.SubmitResultReply{}, s.api.SubmitResult(request.ExecutionID, []byte(request.GetData()), err)
}
