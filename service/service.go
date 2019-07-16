package service

import (
	"fmt"
)

// GetDependency returns dependency dependencyKey or a not found error.
func (s *Service) GetDependency(dependencyKey string) (*Dependency, error) {
	for _, dep := range s.Dependencies {
		if dep.Key == dependencyKey {
			return dep, nil
		}
	}
	return nil, fmt.Errorf("dependency %s do not exist", dependencyKey)
}

// GetTask returns task taskKey of service.
func (s *Service) GetTask(taskKey string) (*Task, error) {
	for _, task := range s.Tasks {
		if task.Key == taskKey {
			return task, nil
		}
	}
	return nil, &TaskNotFoundError{
		TaskKey:     taskKey,
		ServiceName: s.Name,
	}
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

// ValidateTaskInputs produces warnings for task inputs that doesn't satisfy their parameter schemas.
func (s *Service) ValidateTaskInputs(taskKey string, taskInputs map[string]interface{}) ([]*ParameterWarning, error) {
	t, err := s.GetTask(taskKey)
	if err != nil {
		return nil, err
	}
	return validateParametersSchema(t.Inputs, taskInputs), nil
}

// ValidateTaskOutputs produces warnings for task outputs that doesn't satisfy their parameter schemas.
func (s *Service) ValidateTaskOutputs(taskKey string, taskOutputs map[string]interface{}) ([]*ParameterWarning, error) {
	t, err := s.GetTask(taskKey)
	if err != nil {
		return nil, err
	}
	return validateParametersSchema(t.Outputs, taskOutputs), nil
}

// ValidateEventData produces warnings for event datas that doesn't satisfy their parameter schemas.
func (s *Service) ValidateEventData(eventKey string, eventData map[string]interface{}) ([]*ParameterWarning, error) {
	e, err := s.GetEvent(eventKey)
	if err != nil {
		return nil, err
	}
	return validateParametersSchema(e.Data, eventData), nil
}

// RequireTaskInputs requires task inputs to match with parameter schemas.
func (s *Service) RequireTaskInputs(taskKey string, taskInputs map[string]interface{}) error {
	warnings, err := s.ValidateTaskInputs(taskKey, taskInputs)
	if err != nil {
		return err
	}
	if len(warnings) > 0 {
		return &InvalidTaskInputError{
			TaskKey:     taskKey,
			ServiceName: s.Name,
			Warnings:    warnings,
		}
	}
	return nil
}

// RequireTaskOutputs requires task outputs to match with parameter schemas.
func (s *Service) RequireTaskOutputs(taskKey string, taskOutputs map[string]interface{}) error {
	warnings, err := s.ValidateTaskOutputs(taskKey, taskOutputs)
	if err != nil {
		return err
	}
	if len(warnings) > 0 {
		return &InvalidTaskOutputError{
			TaskKey:     taskKey,
			ServiceName: s.Name,
			Warnings:    warnings,
		}
	}
	return nil
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
