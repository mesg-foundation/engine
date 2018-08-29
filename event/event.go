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
	event, err := s.GetEvent(eventKey)
	if err != nil {
		return nil, err
	}
	if err := event.RequireData(eventData); err != nil {
		return nil, err
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
