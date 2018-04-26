package event

import (
	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/types"
)

func getSubscription(request *types.ListenEventRequest) (subscription chan pubsub.Message) {
	service := service.New(request.Service)

	subscription = pubsub.Subscribe(service.EventSubscriptionChannel())
	return
}

// Listen for event from the services
func (s *Server) Listen(request *types.ListenEventRequest, stream types.Event_ListenServer) (err error) {
	subscription := getSubscription(request)
	for data := range subscription {
		stream.Send(data.(*types.EventReply))
	}
	return
}
