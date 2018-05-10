package core

import (
	"encoding/json"

	"github.com/mesg-foundation/core/execution"
	"golang.org/x/net/context"
)

// Execute a task
func (s *Server) ExecuteTask(ctx context.Context, request *ExecuteTaskRequest) (reply *ExecuteTaskReply, err error) {
	service := request.Service
	var inputs interface{}
	err = json.Unmarshal([]byte(request.TaskData), &inputs)
	if err != nil {
		return
	}
	execution, err := execution.Create(service, request.TaskKey, inputs)
	if err != nil {
		return
	}
	err = execution.Execute()
	reply = &ExecuteTaskReply{
		ExecutionID: execution.ID,
	}
	return
}
