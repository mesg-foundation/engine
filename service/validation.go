package service

import (
	"os"
	"path/filepath"

	"github.com/mesg-foundation/core/service/assets"
	"github.com/xeipuuv/gojsonschema"
	yaml "gopkg.in/yaml.v2"
)

// ValidationResult contains the result of the validation of a service
type ValidationResult struct {
	FileWarnings    []string
	DockerfileExist bool
}

// IsValid returns true if all the validation result is valid
func (v *ValidationResult) IsValid() bool {
	if len(v.FileWarnings) == 0 && v.DockerfileExist {
		return true
	}
	return false
}

// ValidateService validates a service at a given path
func ValidateService(path string) (validation ValidationResult, err error) {
	data, err := readFromPath(path)
	if err != nil {
		return
	}
	validation.FileWarnings, err = validateServiceFile(data)
	if err != nil {
		return
	}
	validation.DockerfileExist, err = validateDockerfile(path)
	return
}

func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	}
	return i
}

// validateServiceFile returns a list of warnings (empty if no warning)
// The all validation can be found in https://github.com/mesg-foundation/core/tree/dev/service/schema.json
func validateServiceFile(data []byte) (warnings []string, err error) {
	var body interface{}
	if err = yaml.Unmarshal(data, &body); err != nil {
		return
	}
	body = convert(body)
	schemaData, err := assets.Asset("service/assets/schema.json")
	if err != nil {
		return
	}
	schema := gojsonschema.NewBytesLoader(schemaData)
	loaded := gojsonschema.NewGoLoader(body)
	result, err := gojsonschema.Validate(schema, loaded)
	if err != nil {
		return
	}
	if result.Valid() == false {
		for _, warning := range result.Errors() {
			warnings = append(warnings, warning.String())
		}
	}
	return
}

func validateDockerfile(path string) (exist bool, err error) {
	dockerFile := filepath.Join(path, "Dockerfile")
	file, err := os.Open(dockerFile)
	defer file.Close()
	if os.IsNotExist(err) {
		err = nil
		exist = false
		return
	}
	exist = true
	return
}
