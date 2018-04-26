package task

import (
	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/types"
	"golang.org/x/net/context"
)

// Execute a task
func (s *Server) Execute(ctx context.Context, request *types.ExecuteTaskRequest) (reply *types.TaskReply, err error) {
	channel := service.New(request.Service).TaskSubscriptionChannel()

	reply = &types.TaskReply{
		Task: request.Task,
		Data: request.Data,
	}

	go pubsub.Publish(channel, reply)

	return
}
