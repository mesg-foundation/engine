package event

import (
	"time"

	"github.com/mesg-foundation/core/service"
)

// Create an event
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
			Event: serviceEvent,
			Data:  data,
		}
	}
	return &Event{
		Service:   serviceForEvent,
		Key:       eventKey,
		Data:      data,
		CreatedAt: time.Now(),
	}, nil
}
