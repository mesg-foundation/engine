package event

import (
	"time"

	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/service"
)

// Event stores all informations about Events.
type Event struct {
	Service   *service.Service
	Key       string
	Data      interface{}
	CreatedAt time.Time
}

// Create creates an event.
func Create(s *service.Service, eventKey string, eventData map[string]interface{}) (*Event, error) {
	event, ok := s.Events[eventKey]
	if !ok {
		return nil, &service.EventNotFoundError{
			EventKey:    eventKey,
			ServiceName: s.Name,
		}
	}
	warnings := s.ValidateParametersSchema(event.Data, eventData)
	if len(warnings) > 0 {
		return nil, &service.InvalidEventDataError{
			EventKey:    eventKey,
			ServiceName: s.Name,
			Warnings:    warnings,
		}
	}
	return &Event{
		Service:   s,
		Key:       eventKey,
		Data:      eventData,
		CreatedAt: time.Now(),
	}, nil
}

// Publish publishes an event for every listener.
func (event *Event) Publish() {
	channel := event.Service.EventSubscriptionChannel()
	go pubsub.Publish(channel, event)
}
