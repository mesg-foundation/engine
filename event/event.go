package event

import (
	"time"

	"github.com/mesg-foundation/core/instance"
	"github.com/mesg-foundation/core/service"
)

// Event stores all informations about Events.
type Event struct {
	Instance  *instance.Instance
	Key       string
	Data      interface{}
	CreatedAt time.Time
}

// Create creates an event eventKey with eventData for service s.
func Create(s *service.Service, i *instance.Instance, eventKey string, eventData map[string]interface{}) (*Event, error) {
	if err := s.RequireEventData(eventKey, eventData); err != nil {
		return nil, err
	}
	return &Event{
		Instance:  i,
		Key:       eventKey,
		Data:      eventData,
		CreatedAt: time.Now(),
	}, nil
}
