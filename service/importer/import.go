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
// interface types to their actual definition types if there are.
func normalizeService(def *ServiceDefinition) {
	for _, task := range def.Tasks {
		normalizeParameters(task.Inputs)
		for _, output := range task.Outputs {
			normalizeParameters(output.Data)
		}
	}
	for _, event := range def.Events {
		normalizeParameters(event.Data)
	}
}

func normalizeParameters(params map[string]*Parameter) {
	for _, param := range params {
		normalizeParameter(param)
	}
}

// normalizeParameter normalizes nested parameters by converting
// map[interface{}]interface{} types to map[string]*Parameter types.
func normalizeParameter(param *Parameter) {
	if param == nil {
		return
	}
	if t, ok := param.Type.(map[interface{}]interface{}); ok {
		var params map[string]*Parameter
		mapstructure.Decode(t, &params)
		param.Type = params
		for _, param := range params {
			normalizeParameter(param)
		}
	}
}
