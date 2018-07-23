package core

import (
	"context"
	"encoding/json"

	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/execution"
)

// ExecuteTask will execute a task for a given service
func (s *Server) ExecuteTask(ctx context.Context, request *ExecuteTaskRequest) (reply *ExecuteTaskReply, err error) {
	service, err := services.Get(request.ServiceID)
	if err != nil {
		return
	}
	var inputs interface{}
	err = json.Unmarshal([]byte(request.InputData), &inputs)
	if err != nil {
		return
	}
	exec, err := execution.Create(&service, request.TaskKey, inputs)
	if err != nil {
		return
	}
	err = exec.Execute()
	reply = &ExecuteTaskReply{
		ExecutionID: exec.ID,
	}
	return
}
