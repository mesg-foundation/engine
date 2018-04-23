package event

import (
	"github.com/mesg-foundation/application/pubsub"
	"github.com/mesg-foundation/application/service"
	"github.com/mesg-foundation/application/types"
)

func getSubscription(request *types.ListenEventRequest) (subscription chan pubsub.Message) {
	service := service.New(request.Service)

	subscription = pubsub.Subscribe(service.EventSubscriptionKey())
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
