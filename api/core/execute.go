package core

import (
	"encoding/json"

	"github.com/mesg-foundation/core/database/services"

	"context"

	"github.com/mesg-foundation/core/execution"
)

// ExecuteTask will execute a task for a given service
func (s *Server) ExecuteTask(ctx context.Context, request *ExecuteTaskRequest) (*ExecuteTaskReply, error) {
	service, err := services.Get(request.ServiceID)
	if err != nil {
		return nil, err
	}
	var inputs map[string]interface{}
	err = json.Unmarshal([]byte(request.InputData), &inputs)
	if err != nil {
		return nil, err
	}
	execution, err := execution.Create(&service, request.TaskKey, inputs, request.Tags)
	if err != nil {
		return nil, err
	}
	err = execution.Execute()
	return &ExecuteTaskReply{
		ExecutionID: execution.ID,
	}, err
}
