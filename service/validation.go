package service

import (
	"fmt"
)

// ParameterWarning contains a specific warning related to a parameter.
type ParameterWarning struct {
	Key       string
	Warning   string
	Parameter *Parameter
}

func (p *ParameterWarning) String() string {
	return fmt.Sprintf("Value of %q is %s", p.Key, p.Warning)
}

// ValidateParametersSchema validates data to see if it matches with parameters schema.
// TODO(ilgooz) make this private when we manually create Service type.
// TODO(ilgooz) remove pointer from *Parameter.
func (s *Service) ValidateParametersSchema(parameters map[string]*Parameter,
	data map[string]interface{}) []*ParameterWarning {
	warnings := make([]*ParameterWarning, 0)

	for key, parameter := range parameters {
		warning := newParameterValidator(key, parameter).Validate(data[key])
		if warning != nil {
			warnings = append(warnings, warning)
		}
	}

	return warnings
}

type parameterValidator struct {
	key       string
	parameter *Parameter
}

func newParameterValidator(key string, parameter *Parameter) *parameterValidator {
	return &parameterValidator{key, parameter}
}

// Validate returns a warning based on the match of the data given in parameter and the parameter.
func (v *parameterValidator) Validate(value interface{}) *ParameterWarning {
	if value == nil {
		if v.parameter.Optional {
			return nil
		}
		return v.newParameterWarning("required")
	}

	return v.validateType(value)
}

func (v *parameterValidator) validateType(value interface{}) *ParameterWarning {
	switch v.parameter.Type {
	case "String":
		if _, ok := value.(string); !ok {
			return v.newParameterWarning("not a string")
		}

	case "Number":
		_, okFloat64 := value.(float64)
		_, okFloat32 := value.(float32)
		_, okInt := value.(int)
		if !okInt && !okFloat64 && !okFloat32 {
			return v.newParameterWarning("not a number")
		}

	case "Boolean":
		if _, ok := value.(bool); !ok {
			return v.newParameterWarning("not a boolean")
		}

	case "Object":
		_, okObj := value.(map[string]interface{})
		_, okArr := value.([]interface{})
		if !okObj && !okArr {
			return v.newParameterWarning("not an object or array")
		}

	default:
		return v.newParameterWarning("an invalid type")
	}

	return nil
}

func (v *parameterValidator) newParameterWarning(warning string) *ParameterWarning {
	return &ParameterWarning{
		Key:       v.key,
		Warning:   warning,
		Parameter: v.parameter,
	}
}
