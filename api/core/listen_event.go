package core

import (
	"encoding/json"
	fmt "fmt"

	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/event"
	"github.com/mesg-foundation/core/pubsub"
	service "github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/utils/array"
)

// ListenEvent for listen event from a specific service services
func (s *Server) ListenEvent(request *ListenEventRequest, stream Core_ListenEventServer) error {
	service, err := services.Get(request.ServiceID)
	if err != nil {
		return err
	}
	if err := validateEventKey(&service, request.EventFilter); err != nil {
		return err
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
	return nil
}

func validateEventKey(service *service.Service, eventKey string) error {
	if eventKey == "" || eventKey == "*" {
		return nil
	}
	if _, ok := service.Events[eventKey]; ok {
		return nil
	}
	return fmt.Errorf("Event %q doesn't exist in this service", eventKey)
}

func isSubscribedEvent(request *ListenEventRequest, e *event.Event) bool {
	return array.IncludedIn([]string{"", "*", e.Key}, request.EventFilter)
}
