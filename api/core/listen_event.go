package core

import (
	"encoding/json"
	"errors"

	"github.com/mesg-foundation/core/database/services"
	service "github.com/mesg-foundation/core/service"

	"github.com/mesg-foundation/core/event"
	"github.com/mesg-foundation/core/pubsub"
)

// ListenEvent for listen event from a specific service services
func (s *Server) ListenEvent(request *ListenEventRequest, stream Core_ListenEventServer) (err error) {
	service, err := services.Get(request.ServiceID)
	if err != nil {
		return
	}
	if err = validateEventKey(&service, request.EventKey); err != nil {
		return
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

func validateEventKey(service *service.Service, eventKey string) (err error) {
	if eventKey == "" {
		return
	}
	if eventKey == "*" {
		return
	}
	_, ok := service.Events[eventKey]
	if ok {
		return
	}
	err = errors.New("Invalid eventKey: " + eventKey)
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
