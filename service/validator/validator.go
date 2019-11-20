package validator

import (
	"fmt"
	"strings"

	"github.com/mesg-foundation/engine/service"
	"github.com/mesg-foundation/engine/x/xerrors"
	"github.com/mesg-foundation/engine/x/xvalidator"
	validator "gopkg.in/go-playground/validator.v9"
)

const namespacePrefix = "service."

var validate, translator = xvalidator.NewWithPrefix(namespacePrefix)

// ValidateService validates if service contains proper data.
func ValidateService(s *service.Service) error {
	if err := validateServiceStruct(s); err != nil {
		return err
	}
	return validateServiceData(s)
}

func validateServiceStruct(s *service.Service) error {
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

func validateServiceData(s *service.Service) error {
	var errs xerrors.Errors
	if err := isServiceKeysUnique(s); err != nil {
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
			if !found && depVolumeKey != service.MainServiceKey {
				err := fmt.Errorf("dependencies[%s].volumesFrom is invalid: dependency %q does not exist", dep.Key, depVolumeKey)
				errs = append(errs, err)
			}
		}
	}
	return errs.ErrorOrNil()
}

// isServiceKeysUnique checks uniqueness of service deps/tasks/events/params keys.
func isServiceKeysUnique(s *service.Service) error {
	var errs xerrors.Errors
	exist := make(map[string]bool)
	for _, dep := range s.Dependencies {
		if exist[dep.Key] {
			errs = append(errs, fmt.Errorf("dependencies[%s] already exist", dep.Key))
		}
		exist[dep.Key] = true
	}

	exist = make(map[string]bool)
	for _, task := range s.Tasks {
		if exist[task.Key] {
			errs = append(errs, fmt.Errorf("tasks[%s] already exist", task.Key))
		}
		exist[task.Key] = true
		if err := isServiceParamsUnique(task.Inputs, fmt.Sprintf("tasks[%s].inputs", task.Key)); err != nil {
			errs = append(errs, err)
		}
		if err := isServiceParamsUnique(task.Outputs, fmt.Sprintf("tasks[%s].outputs", task.Key)); err != nil {
			errs = append(errs, err)
		}
	}

	exist = make(map[string]bool)
	for _, event := range s.Events {
		if exist[event.Key] {
			errs = append(errs, fmt.Errorf("events[%s] already exist", event.Key))
		}
		exist[event.Key] = true

		if err := isServiceParamsUnique(event.Data, fmt.Sprintf("events[%s].data", event.Key)); err != nil {
			errs = append(errs, err)
		}
	}
	return errs.ErrorOrNil()
}

// isServiceParamsUnique checks uniqueness of service params.
func isServiceParamsUnique(ps []*service.Service_Parameter, errprefix string) error {
	if len(ps) == 0 {
		return nil
	}

	var errs xerrors.Errors
	existparam := make(map[string]bool)
	for _, p := range ps {
		if existparam[p.Key] {
			errs = append(errs, fmt.Errorf("%s[%s] already exist", errprefix, p.Key))
		}
		existparam[p.Key] = true

		if err := isServiceParamsUnique(p.Object, fmt.Sprintf("%s[%s].object", errprefix, p.Key)); err != nil {
			errs = append(errs, err)
		}
	}
	return errs.ErrorOrNil()
}
