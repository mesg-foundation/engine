package event

import (
	"strings"

	"github.com/mesg-foundation/core/service"
)

type parameterWarning struct {
	parameter string
	warning   string
}

func (p *parameterWarning) String() string {
	return strings.Join([]string{
		p.parameter,
		p.warning,
	}, " ")
}

func exists(service *service.Service, name string) bool {
	for eventName := range service.Events {
		if eventName == name {
			return true
		}
	}
	return false
}

func validParameters(parameters map[string]*service.Parameter, data map[string]interface{}) bool {
	return len(parametersWarnings(parameters, data)) == 0
}

func parametersWarnings(parameters map[string]*service.Parameter, data map[string]interface{}) (warnings []*parameterWarning) {
	warnings = make([]*parameterWarning, 0)
	for key, param := range parameters {
		warning := checkParameterWarning(key, param, data)
		if warning != nil {
			warnings = append(warnings, warning)
		}
	}
	return
}

func checkParameterWarning(key string, parameter *service.Parameter, data map[string]interface{}) (warning *parameterWarning) {
	if data[key] == nil {
		if parameter.Optional {
			return
		}
		warning = &parameterWarning{parameter: key, warning: "required"}
		return
	}
	value := data[key]
	switch parameter.Type {
	case "String":
		_, ok := value.(string)
		if !ok {
			warning = &parameterWarning{parameter: key, warning: "not a string"}
			return
		}
	case "Number":
		_, okFloat64 := value.(float64)
		_, okFloat32 := value.(float32)
		_, okInt := value.(int)
		if !okInt && !okFloat64 && !okFloat32 {
			warning = &parameterWarning{parameter: key, warning: "not a number"}
			return
		}
	case "Boolean":
		_, ok := value.(bool)
		if !ok {
			warning = &parameterWarning{parameter: key, warning: "not a boolean"}
			return
		}
	case "Object":
		_, okObj := value.(map[string]interface{})
		_, okArr := value.([]interface{})
		if !okObj && !okArr {
			warning = &parameterWarning{parameter: key, warning: "not an object/array"}
			return
		}
	}
	return
}
