package event

import (
	"errors"
	"time"

	"github.com/mesg-foundation/core/service"
)

// Create an event
func Create(serviceForEvent *service.Service, eventKey string, data map[string]interface{}) (event *Event, err error) {
	serviceEvent, evenFound := serviceForEvent.Events[eventKey]
	if !evenFound {
		err = errors.New("Event " + eventKey + " doesn't exists in service " + serviceForEvent.Name)
		return
	}
	if serviceEvent.IsValid(data) {
		errorString := "Invalid parameters: "
		for _, warning := range serviceEvent.Validate(data) {
			errorString = errorString + " " + warning.String()
		}
		err = errors.New(errorString)
		return
	}
	event = &Event{
		Service:   serviceForEvent,
		Key:       eventKey,
		Data:      data,
		CreatedAt: time.Now(),
	}
	return
}
