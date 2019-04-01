package api

import (
	"github.com/mesg-foundation/core/event"
)

// EmitEvent emits a MESG event eventKey with eventData for service token.
func (a *API) EmitEvent(token, eventKey string, eventData map[string]interface{}) error {
	s, err := a.db.Get(token)
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
