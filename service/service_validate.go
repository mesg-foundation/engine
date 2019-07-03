package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/mesg-foundation/engine/x/xerrors"
	"github.com/mesg-foundation/engine/x/xvalidator"
	validator "gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

const namespacePrefix = "service."

var validate, translator = newValidator()

// ValidateService validates if service contains proper data.
func ValidateService(service *Service) error {
	if service.Configuration == nil {
		service.Configuration = &Dependency{}
	}
	if err := validateServiceStruct(service); err != nil {
		return err
	}
	return validateServiceData(service)
}

func validateServiceStruct(service *Service) error {
	var errs xerrors.Errors
	// validate service struct based on tag
	if err := validate.Struct(service); err != nil {
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

func validateServiceData(service *Service) error {
	var errs xerrors.Errors
	if err := isServiceKeysUnique(service); err != nil {
		errs = append(errs, err)
	}

	// validate configuration image
	if service.Configuration.Image != "" {
		errs = append(errs, errors.New("configuration.image is not allowed"))
	}

	// validate dependencies image
	for _, dep := range service.Dependencies {
		if dep.Image == "" {
			err := fmt.Errorf("dependencies[%s].image is a required field", dep.Key)
			errs = append(errs, err)
		}
	}

	// validate configuration volumes
	for _, depVolumeKey := range service.Configuration.VolumesFrom {
		found := false
		for _, s := range service.Dependencies {
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
	for _, dep := range service.Dependencies {
		for _, depVolumeKey := range dep.VolumesFrom {
			found := false
			for _, s := range service.Dependencies {
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

// isServiceKeysUnique checks uniqueness of service deps/tasks/events/params keys.
func isServiceKeysUnique(service *Service) error {
	var errs xerrors.Errors
	exist := make(map[string]bool)
	for _, dep := range service.Dependencies {
		if exist[dep.Key] {
			errs = append(errs, fmt.Errorf("dependencies[%s] already exist", dep.Key))
		}
		exist[dep.Key] = true
	}

	exist = make(map[string]bool)
	for _, task := range service.Tasks {
		if exist[task.Key] {
			errs = append(errs, fmt.Errorf("tasks[%s] already exist", task.Key))
		}
		exist[task.Key] = true

		existparam := make(map[string]bool)
		for _, param := range task.Inputs {
			if existparam[param.Key] {
				errs = append(errs, fmt.Errorf("tasks[%s].inputs[%s] already exist", task.Key, param.Key))
			}
			existparam[param.Key] = true
		}

		existparam = make(map[string]bool)
		for _, param := range task.Outputs {
			if existparam[param.Key] {
				errs = append(errs, fmt.Errorf("tasks[%s].outputs[%s] already exist", task.Key, param.Key))
			}
			existparam[param.Key] = true
		}
	}

	exist = make(map[string]bool)
	for _, event := range service.Events {
		if exist[event.Key] {
			errs = append(errs, fmt.Errorf("events[%s] already exist", event.Key))
		}
		exist[event.Key] = true

		existparam := make(map[string]bool)
		for _, param := range event.Data {
			if existparam[param.Key] {
				errs = append(errs, fmt.Errorf("events[%s].data[%s] already exist", event.Key, param.Key))
			}
			existparam[param.Key] = true
		}
	}
	return errs.ErrorOrNil()
}

func newValidator() (*validator.Validate, ut.Translator) {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	validate := validator.New()

	validate.RegisterValidation("env", xvalidator.IsEnv)
	validate.RegisterTranslation("env", trans, func(ut ut.Translator) error {
		return ut.Add("env", "{0} must be a valid env variable name", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("env", fe.Field(), namespacePrefix)
		return t
	})

	validate.RegisterValidation("portmap", xvalidator.IsPortMapping)
	validate.RegisterTranslation("portmap", trans, func(ut ut.Translator) error {
		return ut.Add("portmap", "{0} must be a valid port mapping. Eg: 80 or 80:80", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("portmap", fe.Field(), namespacePrefix)
		return t
	})
	validate.RegisterValidation("domain", xvalidator.IsDomainName)
	validate.RegisterTranslation("domain", trans, func(ut ut.Translator) error {
		return ut.Add("domain", "{0} must respect domain-style notation. Eg: author.name", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("domain", fe.Field())
		return t
	})
	en_translations.RegisterDefaultTranslations(validate, trans)
	return validate, trans
}
