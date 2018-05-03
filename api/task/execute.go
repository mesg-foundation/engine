package task

import (
	"encoding/json"

	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/types"
	"golang.org/x/net/context"
)

// Execute a task
func (s *Server) Execute(ctx context.Context, request *types.ExecuteTaskRequest) (reply *types.TaskReply, err error) {
	service := service.New(request.Service)
	var inputs interface{}
	err = json.Unmarshal([]byte(request.Data), &inputs)
	if err != nil {
		return
	}
	execution, err := execution.Create(service, request.Task, inputs)
	if err != nil {
		return
	}
	reply, err = execution.Execute()
	return
}
