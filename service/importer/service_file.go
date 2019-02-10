package importer

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	validator "gopkg.in/go-playground/validator.v9"
	yaml "gopkg.in/yaml.v2"
)

// ConfigurationDependencyKey is the reserved key of the service's configuration in the dependencies array.
const ConfigurationDependencyKey = "service"

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
	errs := validator.New().Struct(service)
	warnings := []string{}
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			warnings = append(
				warnings,
				fmt.Sprintf("Value %q from %q is not valid.", err.Value(), err.Field()),
			)
		}
	}
	return warnings, nil
}

func validateServiceFileDependencyKey(data interface{}) (warning string) {
	s, _ := data.(map[string]interface{})
	dep, _ := s["dependencies"].(map[string]interface{})
	if dep[ConfigurationDependencyKey] != nil {
		return fmt.Sprintf("cannot use %q as dependency key", ConfigurationDependencyKey)
	}
	return ""
}
