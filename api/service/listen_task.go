package service

import (
	"github.com/mesg-foundation/core/pubsub"
)

func getSubscription(request *ServiceRequest) (subscription chan pubsub.Message) {
	service := request.Service

	subscription = pubsub.Subscribe(service.TaskSubscriptionChannel())
	return
}

// Listen for tasks
func (s *Server) ListenTask(request *ServiceRequest, stream Service_ListenTaskServer) (err error) {
	subscription := getSubscription(request)
	for data := range subscription {
		stream.Send(data.(*TaskData))
	}
	return
}
