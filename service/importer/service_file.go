package importer

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/mesg-foundation/core/x/xvalidator"
	validator "gopkg.in/go-playground/validator.v9"
	yaml "gopkg.in/yaml.v2"
)

func readServiceFile(path string) ([]byte, error) {
	file := filepath.Join(path, "mesg.yml")
	return ioutil.ReadFile(file)
}

// validateServiceFile returns a list of warnings.
func validateServiceFile(data []byte) []string {
	var service ServiceDefinition
	if err := yaml.UnmarshalStrict(data, &service); err != nil {
		return []string{fmt.Sprintf("parse mesg.yml error: %s", err)}
	}
	return validateServiceStruct(&service)
}

func validateServiceStruct(service *ServiceDefinition) []string {
	validate := validator.New()
	validate.RegisterValidation("port", xvalidator.IsPort)
	validate.RegisterValidation("domain", xvalidator.IsDomainName)

	errs := validate.Struct(service)
	warnings := []string{}
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			warnings = append(
				warnings,
				fmt.Sprintf("%s with value %q is invalid: %s %s", err.Namespace(), err.Value(), err.Tag(), err.Param()),
			)
		}
	}

	for k, d := range service.Dependencies {
		if d.Image == "" {
			fmt.Sprintf("dependencies[%s].image is invalid: empty", k)
		}
	}

	return warnings
}
