package service

// Event describes a service task.
type Event struct {
	// Key is the key of event.
	Key string `hash:"name:1" validate:"printascii"`

	// Name is the name of event.
	Name string `hash:"name:2" validate:"printascii"`

	// Description is the description of event.
	Description string `hash:"name:3" validate:"printascii"`

	// Data holds the input parameters of event.
	Data []*Parameter `hash:"name:4" validate:"dive,required"`
}

// GetEvent returns event eventKey of service.
func (s *Service) GetEvent(eventKey string) (*Event, error) {
	for _, event := range s.Events {
		if event.Key == eventKey {
			return event, nil
		}
	}
	return nil, &EventNotFoundError{
		EventKey:    eventKey,
		ServiceName: s.Name,
	}
}

// ValidateEventData produces warnings for event datas that doesn't satisfy their parameter schemas.
func (s *Service) ValidateEventData(eventKey string, eventData map[string]interface{}) ([]*ParameterWarning, error) {
	e, err := s.GetEvent(eventKey)
	if err != nil {
		return nil, err
	}
	return validateParametersSchema(e.Data, eventData), nil
}

// RequireEventData requires event datas to be matched with parameter schemas.
func (s *Service) RequireEventData(eventKey string, eventData map[string]interface{}) error {
	warnings, err := s.ValidateEventData(eventKey, eventData)
	if err != nil {
		return err
	}
	if len(warnings) > 0 {
		return &InvalidEventDataError{
			EventKey:    eventKey,
			ServiceName: s.Name,
			Warnings:    warnings,
		}
	}
	return nil
}
