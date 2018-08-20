package core

import (
	"context"
	"encoding/json"
)

// ExecuteTask executes a task for a given service.
func (s *Server) ExecuteTask(ctx context.Context, request *ExecuteTaskRequest) (*ExecuteTaskReply, error) {
	var inputs map[string]interface{}
	if err := json.Unmarshal([]byte(request.InputData), &inputs); err != nil {
		return nil, err
	}

	executionID, err := s.api.ExecuteTask(request.ServiceID, request.TaskKey, inputs, request.ExecutionTags)
	return &ExecuteTaskReply{
		ExecutionID: executionID,
	}, err
}
