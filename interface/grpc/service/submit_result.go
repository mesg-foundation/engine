package service

import (
	"context"
	"encoding/json"

	"github.com/mesg-foundation/core/execution"
)

// SubmitResult subbmits results of an execution.
func (s *Server) SubmitResult(context context.Context, request *SubmitResultRequest) (*SubmitResultReply, error) {
	execution := execution.InProgress(request.ExecutionID)
	if execution == nil {
		return nil, &MissingExecutionError{
			ID: request.ExecutionID,
		}
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(request.OutputData), &data); err != nil {
		return nil, err
	}
	if err := execution.Complete(request.OutputKey, data); err != nil {
		return nil, err
	}
	return &SubmitResultReply{}, nil
}
