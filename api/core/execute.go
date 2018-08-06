package core

import (
	"context"
	"encoding/json"

	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/execution"
)

// ExecuteTask executes a task for a given service.
func (s *Server) ExecuteTask(ctx context.Context, request *ExecuteTaskRequest) (*ExecuteTaskReply, error) {
	service, err := services.Get(request.ServiceID)
	if err != nil {
		return nil, err
	}

	var inputs map[string]interface{}
	if err := json.Unmarshal([]byte(request.InputData), &inputs); err != nil {
		return nil, err
	}
	execution, err := execution.Create(&service, request.TaskKey, inputs)
	if err != nil {
		return nil, err
	}

	return &ExecuteTaskReply{
		ExecutionID: execution.ID,
	}, execution.Execute()
}
