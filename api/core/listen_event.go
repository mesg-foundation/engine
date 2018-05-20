package core

import (
	"encoding/json"

	"github.com/mesg-foundation/core/database/services"

	"github.com/mesg-foundation/core/event"
	"github.com/mesg-foundation/core/pubsub"
)

// ListenEvent for listen event from a specific service services
func (s *Server) ListenEvent(request *ListenEventRequest, stream Core_ListenEventServer) (err error) {
	service, err := services.Get(request.ServiceID)
	if err != nil {
		return
	}
	subscription := pubsub.Subscribe(service.EventSubscriptionChannel())
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
