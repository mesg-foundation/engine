package event

import (
	"log"

	"github.com/mesg-foundation/application/pubsub"
	"github.com/mesg-foundation/application/service"
	"github.com/mesg-foundation/application/types"
)

// Listen
func (s *Server) Listen(request *types.ListenEventRequest, stream types.Event_ListenServer) (err error) {
	log.Println("Receive listen request", request)

	service := service.New(request.Service)

	subscription := pubsub.Subscribe(service.EventSubscriptionChannel())

	for data := range subscription {
		reply := data.(*types.EventReply)
		stream.Send(reply)
	}

	return
}
