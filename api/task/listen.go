package task

import (
	"github.com/mesg-foundation/application/pubsub"
	"github.com/mesg-foundation/application/service"
	"github.com/mesg-foundation/application/types"
)

func getSubscription(request *types.ListenTaskRequest) (subscription chan pubsub.Message) {
	service := service.New(request.Service)

	subscription = pubsub.Subscribe(service.TaskSubscriptionKey())
	return
}

// Listen for tasks
func (s *Server) Listen(request *types.ListenTaskRequest, stream types.Task_ListenServer) (err error) {
	subscription := getSubscription(request)
	for data := range subscription {
		stream.Send(data.(*types.TaskReply))
	}
	return
}
