package core

import (
	"encoding/json"
	"errors"

	"github.com/mesg-foundation/core/database/services"
	service "github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/utils/array"

	"github.com/mesg-foundation/core/event"
	"github.com/mesg-foundation/core/pubsub"
)

// ListenEvent for listen event from a specific service services
func (s *Server) ListenEvent(request *ListenEventRequest, stream Core_ListenEventServer) (err error) {
	service, err := services.Get(request.ServiceID)
	if err != nil {
		return
	}
	if err = validateEventKey(&service, request.EventFilter); err != nil {
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
	if eventKey == "" || eventKey == "*" {
		return
	}
	_, ok := service.Events[eventKey]
	if ok {
		return
	}
	err = errors.New("Event '" + eventKey + "' doesn't exist in this service")
	return
}

func isSubscribedEvent(request *ListenEventRequest, e *event.Event) bool {
	return array.IncludedIn([]string{"", "*", e.Key}, request.EventFilter)
}
