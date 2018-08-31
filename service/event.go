package service

// Event describes a service task.
type Event struct {
	// Key is the key of event.
	Key string `hash:"-"`

	// Name is the name of event.
	Name string `hash:"name:1"`

	// Description is the description of event.
	Description string `hash:"name:2"`

	// Data holds the input parameters of event.
	Data []*Parameter `hash:"name:3"`

	// service is the event's service.
	service *Service `hash:"-"`
}

// GetEvent returns event eventKey of service.
func (s *Service) GetEvent(eventKey string) (*Event, error) {
	for _, event := range s.Events {
		if event.Key == eventKey {
			event.service = s
			return event, nil
		}
	}
	return nil, &EventNotFoundError{
		EventKey:    eventKey,
		ServiceName: s.Name,
	}
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
			ServiceName: e.service.Name,
			Warnings:    warnings,
		}
	}
	return nil
}
