package service

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/xeipuuv/gojsonschema"
	yaml "gopkg.in/yaml.v2"
)

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

// validServiceFileData returns a list of warnings (empty if no warning)
// The all validation can be found in https://github.com/mesg-foundation/core/tree/dev/service/schema.json
func validServiceData(data []byte) (warnings []gojsonschema.ResultError, err error) {
	var body interface{}
	if err = yaml.Unmarshal(data, &body); err != nil {
		return
	}
	body = convert(body)
	schemaData, err := Asset("service/schema.json")
	if err != nil {
		return
	}
	schema := gojsonschema.NewBytesLoader(schemaData)
	loaded := gojsonschema.NewGoLoader(body)
	result, err := gojsonschema.Validate(schema, loaded)
	if err != nil {
		return
	}

	if !result.Valid() {
		warnings = result.Errors()
		err = errors.New("File 'mesg.yml' is not valid")
	}
	return
}

// ValidService validates a service at a given path
func ValidService(path string) (warnings []gojsonschema.ResultError, err error) {
	data, err := readFromPath(path)
	if err != nil {
		return
	}

	warnings, err = validServiceData(data)
	if err != nil {
		return
	}

	dockerFile := filepath.Join(path, "Dockerfile")
	file, err := os.Open(dockerFile)
	defer file.Close()
	return
}
