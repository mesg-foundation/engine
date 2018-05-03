package event

import (
	"errors"
	"time"

	"github.com/mesg-foundation/core/service"
)

// Create an event
func Create(service *service.Service, eventKey string, data interface{}) (event *Event, err error) {
	if !eventExists(service, eventKey) {
		err = errors.New("Event " + eventKey + " doesn't exists in service " + service.Name)
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

func eventExists(service *service.Service, name string) (exists bool) {
	exists = false
	for eventName := range service.Events {
		if eventName == name {
			exists = true
			break
		}
	}
	return
}
