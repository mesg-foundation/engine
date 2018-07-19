package event

import (
	"errors"
	"time"

	"github.com/mesg-foundation/core/service"
)

// Create an event
func Create(service *service.Service, eventKey string, data map[string]interface{}) (event *Event, err error) {
	if !exists(service, eventKey) {
		err = errors.New("Event " + eventKey + " doesn't exists in service " + service.Name)
		return
	}
	parameters := service.Events[eventKey].Data
	if !validParameters(parameters, data) {
		errorString := "Invalid parameters: "
		for _, warning := range parametersWarnings(parameters, data) {
			errorString = errorString + " " + warning.String()
		}
		err = errors.New(errorString)
		return
	}
	event = &Event{
		Service:   service,
		Key:       eventKey,
		Data:      data,
		CreatedAt: time.Now(),
	}
	return
}
