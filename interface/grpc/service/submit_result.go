package service

import (
	"context"
	"encoding/json"

	"github.com/mesg-foundation/core/protobuf/serviceapi"
)

// SubmitResult submits results of an execution.
func (s *Server) SubmitResult(context context.Context, request *serviceapi.SubmitResultRequest) (*serviceapi.SubmitResultReply, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(request.OutputData), &data); err != nil {
		return nil, err
	}
	return &serviceapi.SubmitResultReply{}, s.api.SubmitResult(request.ExecutionID, request.OutputKey, data)
}
