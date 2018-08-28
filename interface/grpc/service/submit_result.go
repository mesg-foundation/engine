package service

import (
	"context"
	"encoding/json"
)

// SubmitResult submits results of an execution.
func (s *Server) SubmitResult(context context.Context, request *SubmitResultRequest) (*SubmitResultReply, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(request.OutputData), &data); err != nil {
		return nil, err
	}
	return &SubmitResultReply{}, s.api.SubmitResult(request.ExecutionID, request.OutputKey, data)
}
