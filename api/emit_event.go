package api

import (
	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/event"
)

// EmitEvent emits a MESG event eventKey with eventData for service token.
func (a *API) EmitEvent(token, eventKey string, eventData map[string]interface{}) error {
	return newEventEmitter(a).Emit(token, eventKey, eventData)
}

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
	event, err := event.Create(s, eventKey, eventData)
	if err != nil {
		return err
	}
	event.Publish()
	return nil
}
