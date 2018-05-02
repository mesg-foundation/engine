package client

import (
	"github.com/mesg-foundation/core/pubsub"
)

// Listen for event from the services
func (s *Server) ListenEvent(request *ServiceRequest, stream Client_ListenEventServer) (err error) {
	subscription := pubsub.Subscribe(request.Service.EventSubscriptionChannel())
	for data := range subscription {
		stream.Send(data.(*EventData))
	}
	return
}
