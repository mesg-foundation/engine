package api

// EmitEvent emits a MESG event eventKey with eventData for service token.
func (a *API) EmitEvent(token, eventKey string, eventData map[string]interface{}) error {
	return newEventEmitter(a).Emit(token, eventKey, eventData)
}
