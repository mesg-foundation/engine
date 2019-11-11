package service

import (
	"fmt"

	"github.com/mesg-foundation/engine/protobuf/types"
)

// MainServiceKey is key for main service.
const MainServiceKey = "service"

// GetDependency returns dependency dependencyKey or a not found error.
func (s *Service) GetDependency(dependencyKey string) (*Service_Dependency, error) {
	for _, dep := range s.Dependencies {
		if dep.Key == dependencyKey {
			return dep, nil
		}
	}
	return nil, fmt.Errorf("service %q - dependency %s does not exist", s.Name, dependencyKey)
}

// GetTask returns task taskKey of service.
func (s *Service) GetTask(taskKey string) (*Service_Task, error) {
	for _, task := range s.Tasks {
		if task.Key == taskKey {
			return task, nil
		}
	}
	return nil, fmt.Errorf("service %q - task %q not found", s.Name, taskKey)
}

// GetEvent returns event eventKey of service.
func (s *Service) GetEvent(eventKey string) (*Service_Event, error) {
	for _, event := range s.Events {
		if event.Key == eventKey {
			return event, nil
		}
	}
	return nil, fmt.Errorf("service %q - event %q not found", s.Name, eventKey)
}

// RequireTaskInputs requires task inputs to match with parameter schemas.
func (s *Service) RequireTaskInputs(taskKey string, inputs []*types.Value) error {
	t, err := s.GetTask(taskKey)
	if err != nil {
		return err
	}
	if err := validateServiceParameters(t.Inputs, inputs); err != nil {
		return fmt.Errorf("service %q - inputs of task %q are invalid: %s", s.Name, taskKey, err)
	}
	return nil
}

// RequireTaskOutputs requires task outputs to match with parameter schemas.
func (s *Service) RequireTaskOutputs(taskKey string, outputs []*types.Value) error {
	t, err := s.GetTask(taskKey)
	if err != nil {
		return err
	}
	if err := validateServiceParameters(t.Outputs, outputs); err != nil {
		return fmt.Errorf("service %q - outputs of task %q are invalid: %s", s.Name, taskKey, err)
	}
	return nil
}

// RequireEventData requires event datas to be matched with parameter schemas.
func (s *Service) RequireEventData(eventKey string, data []*types.Value) error {
	e, err := s.GetEvent(eventKey)
	if err != nil {
		return err
	}
	if err := validateServiceParameters(e.Data, data); err != nil {
		return fmt.Errorf("service %q - data of event %q are invalid: %s", s.Name, eventKey, err)
	}
	return nil
}
