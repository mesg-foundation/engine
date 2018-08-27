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
func Create(serviceForEvent *service.Service, eventKey string, data map[string]interface{}) (*Event, error) {
	serviceEvent, eventFound := serviceForEvent.Events[eventKey]
	if !eventFound {
		return nil, &service.EventNotFoundError{
			EventKey:    eventKey,
			ServiceName: serviceForEvent.Name,
		}
	}
	warnings := serviceForEvent.ValidateParametersSchema(serviceEvent.Data, data)
	if len(warnings) > 0 {
		return nil, &service.InvalidEventDataError{
			EventKey: eventKey,
			Warnings: warnings,
		}
	}
	return &Event{
		Service:   serviceForEvent,
		Key:       eventKey,
		Data:      data,
		CreatedAt: time.Now(),
	}, nil
}

// Publish publishes an event for every listener.
func (event *Event) Publish() {
	channel := event.Service.EventSubscriptionChannel()
	go pubsub.Publish(channel, event)
}
