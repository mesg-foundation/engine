package importer

import (
	"github.com/mitchellh/mapstructure"
	yaml "gopkg.in/yaml.v2"
)

// From imports a service from a source.
func From(source string) (*ServiceDefinition, error) {
	return fromPath(source)
}

// fromPath imports a service from a path.
func fromPath(path string) (*ServiceDefinition, error) {
	isValid, err := IsValid(path)
	if err != nil {
		return nil, err
	}
	if !isValid {
		return nil, &ValidationError{}
	}
	data, err := readServiceFile(path)
	if err != nil {
		return nil, err
	}
	var def *ServiceDefinition
	if err := yaml.UnmarshalStrict(data, &def); err != nil {
		return nil, err
	}
	normalizeService(def)
	return def, nil
}

// normalizeService normalizes types in service definition by converting
// interface types to their actual definition type if there is.
func normalizeService(def *ServiceDefinition) {
	for _, task := range def.Tasks {
		for _, input := range task.Inputs {
			normalizeParameter(input)
		}
		for _, output := range task.Outputs {
			for _, d := range output.Data {
				normalizeParameter(d)
			}
		}
	}
	for _, event := range def.Events {
		for _, d := range event.Data {
			normalizeParameter(d)
		}
	}
}

// normalizeParameter normalizes nested parameters to have map[string]*Parameter type
// instead of map[interface{}]interface{}.
func normalizeParameter(param *Parameter) {
	if param == nil {
		return
	}
	switch t := param.Type.(type) {
	case map[interface{}]interface{}:
		var params map[string]*Parameter
		mapstructure.Decode(t, &params)
		param.Type = params
		for _, param := range params {
			normalizeParameter(param)
		}
	}
}
