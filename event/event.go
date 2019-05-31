package event

import (
	"time"

	"github.com/mesg-foundation/core/service"
)

// Event stores all informations about Events.
type Event struct {
	Service   *service.Service
	Key       string
	Data      interface{}
	CreatedAt time.Time
}

// Create creates an event eventKey with eventData for service s.
func Create(s *service.Service, eventKey string, eventData map[string]interface{}) (*Event, error) {
	if err := s.RequireEventData(eventKey, eventData); err != nil {
		return nil, err
	}
	return &Event{
		Service:   s,
		Key:       eventKey,
		Data:      eventData,
		CreatedAt: time.Now(),
	}, nil
}
