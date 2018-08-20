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
			Service:  serviceForEvent,
			EventKey: eventKey,
		}
	}
	if !serviceEvent.IsValid(data) {
		return nil, &service.InvalidEventDataError{
			Event:     serviceEvent,
			EventKey:  eventKey,
			EventData: data,
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
