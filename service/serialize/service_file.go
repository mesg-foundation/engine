package serialize

import (
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/mesg-foundation/core/service/serialize/assets"
	"github.com/xeipuuv/gojsonschema"
	yaml "gopkg.in/yaml.v2"
)

func readServiceFile(path string) (data []byte, err error) {
	file := filepath.Join(path, "mesg.yml")
	data, err = ioutil.ReadFile(file)
	return
}

// validateServiceFile returns a list of warnings (empty if no warning)
func validateServiceFile(data []byte) (warnings []string, err error) {
	var body interface{}
	err = yaml.Unmarshal(data, &body)
	if err != nil {
		err = errors.New("Error with file 'mesg.yml'. " + err.Error())
		return
	}
	body = convert(body)
	result, err := validateServiceFileSchema(body, "service/serialize/assets/schema.json")
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

func validateServiceFileSchema(data interface{}, schemaPath string) (result *gojsonschema.Result, err error) {
	schemaData, err := assets.Asset(schemaPath)
	if err != nil {
		return
	}
	schema := gojsonschema.NewBytesLoader(schemaData)
	loaded := gojsonschema.NewGoLoader(data)
	result, err = gojsonschema.Validate(schema, loaded)
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
