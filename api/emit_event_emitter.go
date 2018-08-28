package api

import (
	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/event"
	"github.com/mesg-foundation/core/service"
)

// eventEmitter provides functionalities to emit a MESG event.
type eventEmitter struct {
	api *API
}

// newEventEmitter creates a new eventEmitter with given api.
func newEventEmitter(api *API) *eventEmitter {
	return &eventEmitter{
		api: api,
	}
}

// Emit emits a MESG event eventKey with eventData for service token.
func (e *eventEmitter) Emit(token, eventKey string, eventData map[string]interface{}) error {
	s, err := services.Get(token)
	if err != nil {
		return err
	}

	serviceEvent, eventFound := s.Events[eventKey]
	if !eventFound {
		return &service.EventNotFoundError{
			EventKey:    eventKey,
			ServiceName: s.Name,
		}
	}
	warnings := s.ValidateParametersSchema(serviceEvent.Data, eventData)
	if len(warnings) > 0 {
		return &service.InvalidEventDataError{
			EventKey: eventKey,
			Warnings: warnings,
		}
	}

	event, err := event.Create(&s, eventKey, eventData)
	if err != nil {
		return err
	}
	event.Publish()
	return nil
}
