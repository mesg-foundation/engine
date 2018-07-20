package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/mesg-foundation/core/execution"
)

// SubmitResult of an execution
func (s *Server) SubmitResult(context context.Context, request *SubmitResultRequest) (*SubmitResultReply, error) {
	execution := execution.InProgress(request.ExecutionID)
	if execution == nil {
		return nil, errors.New("No task in progress with the ID " + request.ExecutionID)
	}
	var data map[string]interface{}
	err := json.Unmarshal([]byte(request.OutputData), &data)
	if err != nil {
		return nil, err
	}
	err = execution.Complete(request.OutputKey, data)
	if err != nil {
		return nil, err
	}
	return &SubmitResultReply{}, nil
}
