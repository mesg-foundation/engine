package service

import (
	"strings"
)

// EventNotFoundError is an error when we cannot find an event in a service
type EventNotFoundError struct {
	Service  *Service
	EventKey string
}

func (e *EventNotFoundError) Error() string {
	return strings.Join([]string{
		"Event",
		e.EventKey,
		"not found in service",
		e.Service.Name,
	}, " ")
}

// InvalidEventDataError is an error when the data of an event are not valid
type InvalidEventDataError struct {
	Event *Event
	Key   string
	Data  map[string]interface{}
}

func (e *InvalidEventDataError) Error() string {
	errorString := "Parameter " + e.Key + " is"
	for _, warning := range e.Event.Validate(e.Data) {
		errorString = errorString + " " + warning.String()
	}
	return errorString
}
