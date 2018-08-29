package service

// GetEvent returns event eventKey of service.
func (s *Service) GetEvent(eventKey string) (*Event, error) {
	event, ok := s.Events[eventKey]
	if !ok {
		return nil, &EventNotFoundError{
			EventKey:    eventKey,
			ServiceName: s.Name,
		}
	}
	return event, nil
}

// ValidateData produces warnings for event datas that doesn't satisfy their parameter schemas.
func (e *Event) ValidateData(eventData map[string]interface{}) []*ParameterWarning {
	return validateParametersSchema(e.Data, eventData)
}

// RequireData requires event datas to be matched with parameter schemas.
func (e *Event) RequireData(eventData map[string]interface{}) error {
	warnings := e.ValidateData(eventData)
	if len(warnings) > 0 {
		return &InvalidEventDataError{
			EventKey:    e.Key,
			ServiceName: e.ServiceName,
			Warnings:    warnings,
		}
	}
	return nil
}
