package importer

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	validator "gopkg.in/go-playground/validator.v9"
	yaml "gopkg.in/yaml.v2"
)

func readServiceFile(path string) ([]byte, error) {
	file := filepath.Join(path, "mesg.yml")
	return ioutil.ReadFile(file)
}

// validateServiceFile returns a list of warnings.
func validateServiceFile(data []byte) ([]string, error) {
	var service ServiceDefinition
	if err := yaml.UnmarshalStrict(data, &service); err != nil {
		return nil, fmt.Errorf("parse mesg.yml error: %s", err)
	}
	return validateServiceStruct(&service)
}

func validateServiceStruct(service *ServiceDefinition) ([]string, error) {
	errs := validator.New().Struct(service)
	warnings := []string{}
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			warnings = append(
				warnings,
				fmt.Sprintf("%s with value %q is invalid: %s %s", err.Field(), err.Value(), err.ActualTag(), err.Param()),
			)
		}
	}
	return warnings, nil
}
