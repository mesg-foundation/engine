package core

import (
	"encoding/json"
	"errors"

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
	if request.EventKey != "" && request.EventKey != "*" {
		_, ok := service.Events[request.EventKey]
		if !ok {
			err = errors.New("Invalid eventKey: " + request.EventKey)
			return
		}
	}
	subscription := pubsub.Subscribe(service.EventSubscriptionChannel())
	for data := range subscription {
		event := data.(*event.Event)
		if isSubscribedEvent(request, event) {
			eventData, _ := json.Marshal(event.Data)
			stream.Send(&EventData{
				EventKey:  event.Key,
				EventData: string(eventData),
			})
		}
	}
	return
}

func isSubscribedEvent(request *ListenEventRequest, e *event.Event) bool {
	if request.EventKey != "" && request.EventKey != "*" && request.EventKey != e.Key {
		return false
	}
	// Possibility to add more filters here like filters on data, awlays return the
	// falsy value and go until the end to have the truth value
	return true
}
