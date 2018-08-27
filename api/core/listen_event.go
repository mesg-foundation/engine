package core

import (
	"encoding/json"
	"fmt"

	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/event"
	"github.com/mesg-foundation/core/pubsub"
	service "github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/x/xstrings"
)

// ListenEvent listens for an event from a specific service.
func (s *Server) ListenEvent(request *ListenEventRequest, stream Core_ListenEventServer) error {
	service, err := services.Get(request.ServiceID)
	if err != nil {
		return err
	}
	if err := validateEventKey(&service, request.EventFilter); err != nil {
		return err
	}

	ctx := stream.Context()
	channel := service.EventSubscriptionChannel()
	subscription := pubsub.Subscribe(channel)
	defer pubsub.Unsubscribe(channel, subscription)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data := <-subscription:
			event := data.(*event.Event)
			if isSubscribedEvent(request, event) {
				eventData, _ := json.Marshal(event.Data)
				if err := stream.Send(&EventData{
					EventKey:  event.Key,
					EventData: string(eventData),
				}); err != nil {
					return err
				}
			}
		}
	}
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
	return xstrings.SliceContains([]string{"", "*", e.Key}, request.EventFilter)
}
