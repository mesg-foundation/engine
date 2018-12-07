package importer

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/mesg-foundation/core/service/importer/assets"
	"github.com/xeipuuv/gojsonschema"
	yaml "gopkg.in/yaml.v2"
)

func readServiceFile(path string) ([]byte, error) {
	file := filepath.Join(path, "mesg.yml")
	return ioutil.ReadFile(file)
}

// validateServiceFile returns a list of warnings.
func validateServiceFile(data []byte) ([]string, error) {
	var body interface{}
	if err := yaml.Unmarshal(data, &body); err != nil {
		return nil, fmt.Errorf("parse mesg.yml error: %s", err)
	}
	body = convert(body)
	result, err := validateServiceFileSchema(body, "service/importer/assets/schema.json")
	if err != nil {
		return nil, err
	}
	var warnings []string
	if !result.Valid() {
		for _, warning := range result.Errors() {
			warnings = append(warnings, warning.String())
		}
	}
	return warnings, nil
}

func validateServiceFileSchema(data interface{}, schemaPath string) (*gojsonschema.Result, error) {
	schemaData, err := assets.Asset(schemaPath)
	if err != nil {
		return nil, err
	}
	schema := gojsonschema.NewBytesLoader(schemaData)
	loaded := gojsonschema.NewGoLoader(data)
	return gojsonschema.Validate(schema, loaded)
}

func convert(i interface{}) interface{} {
	x, ok := i.(map[interface{}]interface{})
	if !ok {
		return i
	}
	m2 := map[string]interface{}{}
	for k, v := range x {
		m2[k.(string)] = convert(v)
	}
	return m2
}
