package service

import (
	"fmt"
	"strings"

	validator "github.com/go-playground/validator/v10"
	"github.com/mesg-foundation/engine/ext/xerrors"
	"github.com/mesg-foundation/engine/ext/xvalidator"
)

const namespacePrefix = "service."

var validate, translator = xvalidator.New(namespacePrefix)

// Validate validates if service contains proper data.
func (s *Service) Validate() error {
	if err := s.validateStruct(); err != nil {
		return err
	}
	return s.validateData()
}

func (s *Service) validateStruct() error {
	var errs xerrors.Errors
	// validate service struct based on tag
	if err := validate.Struct(s); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			// Remove the name of the struct and the field from namespace
			trimmedNamespace := strings.TrimPrefix(e.Namespace(), namespacePrefix)
			trimmedNamespace = strings.TrimSuffix(trimmedNamespace, e.Field())
			// Only use it when in-cascade field
			namespace := ""
			if e.Field() != trimmedNamespace {
				namespace = trimmedNamespace
			}
			errs = append(errs, fmt.Errorf("%s%s", namespace, e.Translate(translator)))
		}
	}
	return errs.ErrorOrNil()
}

func (s *Service) validateData() error {
	var errs xerrors.Errors
	if err := s.areKeysUnique(); err != nil {
		errs = append(errs, err)
	}

	// validate dependencies image
	for _, dep := range s.Dependencies {
		if dep.Image == "" {
			err := fmt.Errorf("dependencies[%s].image is a required field", dep.Key)
			errs = append(errs, err)
		}
	}

	// validate configuration volumes
	for _, depVolumeKey := range s.Configuration.VolumesFrom {
		found := false
		for _, s := range s.Dependencies {
			if s.Key == depVolumeKey {
				found = true
				break
			}
		}
		if !found {
			err := fmt.Errorf("configuration.volumesFrom is invalid: dependency %q does not exist", depVolumeKey)
			errs = append(errs, err)
		}
	}

	// validate dependencies volumes
	for _, dep := range s.Dependencies {
		for _, depVolumeKey := range dep.VolumesFrom {
			found := false
			for _, s := range s.Dependencies {
				if s.Key == depVolumeKey {
					found = true
					break
				}
			}
			if !found && depVolumeKey != MainServiceKey {
				err := fmt.Errorf("dependencies[%s].volumesFrom is invalid: dependency %q does not exist", dep.Key, depVolumeKey)
				errs = append(errs, err)
			}
		}
	}
	return errs.ErrorOrNil()
}

// areKeysUnique checks uniqueness of service deps/tasks/events/params keys.
func (s *Service) areKeysUnique() error {
	var errs xerrors.Errors
	exist := make(map[string]bool)
	for _, dep := range s.Dependencies {
		if exist[dep.Key] {
			errs = append(errs, fmt.Errorf("dependencies[%s] already exists", dep.Key))
		}
		exist[dep.Key] = true
	}

	exist = make(map[string]bool)
	for _, task := range s.Tasks {
		if exist[task.Key] {
			errs = append(errs, fmt.Errorf("tasks[%s] already exists", task.Key))
		}
		exist[task.Key] = true
		if err := areServiceParamsUnique(task.Inputs, fmt.Sprintf("tasks[%s].inputs", task.Key)); err != nil {
			errs = append(errs, err)
		}
		if err := areServiceParamsUnique(task.Outputs, fmt.Sprintf("tasks[%s].outputs", task.Key)); err != nil {
			errs = append(errs, err)
		}
	}

	exist = make(map[string]bool)
	for _, event := range s.Events {
		if exist[event.Key] {
			errs = append(errs, fmt.Errorf("events[%s] already exists", event.Key))
		}
		exist[event.Key] = true

		if err := areServiceParamsUnique(event.Data, fmt.Sprintf("events[%s].data", event.Key)); err != nil {
			errs = append(errs, err)
		}
	}
	return errs.ErrorOrNil()
}

// areServiceParamsUnique checks uniqueness of service params.
func areServiceParamsUnique(ps []*Service_Parameter, errprefix string) error {
	if len(ps) == 0 {
		return nil
	}

	var errs xerrors.Errors
	existparam := make(map[string]bool)
	for _, p := range ps {
		if existparam[p.Key] {
			errs = append(errs, fmt.Errorf("%s[%s] already exists", errprefix, p.Key))
		}
		existparam[p.Key] = true

		if err := areServiceParamsUnique(p.Object, fmt.Sprintf("%s[%s].object", errprefix, p.Key)); err != nil {
			errs = append(errs, err)
		}
	}
	return errs.ErrorOrNil()
}
