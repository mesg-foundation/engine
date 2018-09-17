package service

import (
	"context"
	"encoding/json"

	"github.com/mesg-foundation/core/protobuf/service"
)

// SubmitResult submits results of an execution.
func (s *Server) SubmitResult(context context.Context, request *service.SubmitResultRequest) (*service.SubmitResultReply, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(request.OutputData), &data); err != nil {
		return nil, err
	}
	return &service.SubmitResultReply{}, s.api.SubmitResult(request.ExecutionID, request.OutputKey, data)
}
