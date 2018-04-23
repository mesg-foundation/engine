package task

import (
	"github.com/mesg-foundation/application/pubsub"
	"github.com/mesg-foundation/application/service"
	"github.com/mesg-foundation/application/types"
	"github.com/mesg-foundation/application/utils/hash"
)

// Listen for tasks
func (s *Server) Listen(request *types.ListenTaskRequest, stream types.Task_ListenServer) (err error) {
	service := service.New(request.Service)

	subscription := pubsub.Subscribe(hash.Calculate([]string{
		service.Name,
		"Task",
	}))

	for data := range subscription {
		taskReply := data.(*types.TaskReply)
		stream.Send(taskReply)
	}

	return
}
