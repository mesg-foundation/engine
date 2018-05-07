package client

import (
	"encoding/json"

	"github.com/mesg-foundation/core/event"
	"github.com/mesg-foundation/core/pubsub"
)

// Listen for event from the services
func (s *Server) ListenEvent(request *ListenEventRequest, stream Core_ListenEventServer) (err error) {
	subscription := pubsub.Subscribe(request.Service.EventSubscriptionChannel())
	for data := range subscription {
		event := data.(*event.Event)
		eventData, _ := json.Marshal(event.Data)
		stream.Send(&EventData{
			EventKey:  event.Key,
			EventData: string(eventData),
		})
	}
	return
}
