package event

import (
	"log"

	"github.com/mesg-foundation/application/pubsub"
	"github.com/mesg-foundation/application/service"
	"github.com/mesg-foundation/application/types"
)

func getSubscription(request *types.ListenEventRequest) (subscription chan pubsub.Message) {
	service := service.New(request.Service)

	subscription = pubsub.Subscribe(service.EventSubscriptionChannel())
	return
}

// Listen for event from the services
func (s *Server) Listen(request *types.ListenEventRequest, stream types.Event_ListenServer) (err error) {
	log.Println("Receive listen request", request)
	subscription := getSubscription(request)
	for data := range subscription {
		stream.Send(data.(*types.EventReply))
	}
	return
}
