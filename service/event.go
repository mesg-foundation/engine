// Copyright 2018 MESG Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

// Event describes a service task.
type Event struct {
	// Key is the key of event.
	Key string `hash:"name:1"`

	// Name is the name of event.
	Name string `hash:"name:2"`

	// Description is the description of event.
	Description string `hash:"name:3"`

	// Data holds the input parameters of event.
	Data []*Parameter `hash:"name:4"`

	// serviceName is the event's service's name.
	serviceName string `hash:"-"`
}

// GetEvent returns event eventKey of service.
func (s *Service) GetEvent(eventKey string) (*Event, error) {
	for _, event := range s.Events {
		if event.Key == eventKey {
			event.serviceName = s.Name
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
			ServiceName: e.serviceName,
			Warnings:    warnings,
		}
	}
	return nil
}
