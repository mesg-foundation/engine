package event

import (
	"errors"
	"time"

	"github.com/mesg-foundation/core/service"
)

// Create an event
func Create(serviceForEvent *service.Service, eventKey string, data map[string]interface{}) (*Event, error) {
	serviceEvent, evenFound := serviceForEvent.Events[eventKey]
	if !evenFound {
		return nil, errors.New("Event " + eventKey + " doesn't exists in service " + serviceForEvent.Name)
	}
	if !serviceEvent.IsValid(data) {
		errorString := "Invalid parameters: "
		for _, warning := range serviceEvent.Validate(data) {
			errorString = errorString + " " + warning.String()
		}
		return nil, errors.New(errorString)
	}
	return &Event{
		Service:   serviceForEvent,
		Key:       eventKey,
		Data:      data,
		CreatedAt: time.Now(),
	}, nil
}
