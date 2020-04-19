package service

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/tendermint/tendermint/crypto"
)

// New initializes a new Service.
func New(sid, name, description string, configuration Service_Configuration, tasks []*Service_Task, events []*Service_Event, dependencies []*Service_Dependency, repository, source string) (*Service, error) {
	// create service
	srv := &Service{
		Sid:           sid,
		Name:          name,
		Description:   description,
		Configuration: configuration,
		Tasks:         tasks,
		Events:        events,
		Dependencies:  dependencies,
		Repository:    repository,
		Source:        source,
	}

	// calculate and apply hash to service.
	srv.Hash = hash.Dump(srv)
	srv.Address = sdk.AccAddress(crypto.AddressHash(srv.Hash))

	// set a sid if this one is empty (yes, after hash calculation..)
	if srv.Sid == "" {
		// make sure that sid doesn't have the same length with id.
		srv.Sid = "_" + srv.Hash.String()
	}

	return srv, srv.Validate()
}

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
func (s *Service) RequireTaskInputs(taskKey string, inputs *types.Struct) error {
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
func (s *Service) RequireTaskOutputs(taskKey string, outputs *types.Struct) error {
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
func (s *Service) RequireEventData(eventKey string, data *types.Struct) error {
	e, err := s.GetEvent(eventKey)
	if err != nil {
		return err
	}
	if err := validateServiceParameters(e.Data, data); err != nil {
		return fmt.Errorf("service %q - data of event %q are invalid: %s", s.Name, eventKey, err)
	}
	return nil
}
