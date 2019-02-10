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

// New creates an event eventKey with eventData for given service.
func New(s *service.Service, eventKey string, eventData map[string]interface{}) *Event {
	return &Event{
		Service:   s,
		Key:       eventKey,
		Data:      eventData,
		CreatedAt: time.Now(),
	}
}
