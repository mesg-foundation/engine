package task

import (
	"github.com/mesg-foundation/application/pubsub"
	"github.com/mesg-foundation/application/service"
	"github.com/mesg-foundation/application/types"
	"github.com/mesg-foundation/application/utils/hash"
	"golang.org/x/net/context"
)

// Execute a task
func (s *Server) Execute(ctx context.Context, request *types.ExecuteTaskRequest) (reply *types.TaskReply, err error) {
	service := service.New(request.Service)

	reply = &types.TaskReply{
		Data: request.Data,
		Task: request.Task,
	}

	go pubsub.Publish(hash.Calculate([]string{
		service.Name,
		"Task",
	}), reply)

	return
}
