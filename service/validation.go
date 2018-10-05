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
// TODO(ilgooz) add this as a method to Service type when create custom types for Event, Task etc.
// TODO(ilgooz) remove pointer from *Parameter.
func validateParametersSchema(parameters []*Parameter,
	data map[string]interface{}) []*ParameterWarning {
	warnings := make([]*ParameterWarning, 0)

	for _, param := range parameters {
		warning := newParameterValidator(param).Validate(data[param.Key])
		if warning != nil {
			warnings = append(warnings, warning)
		}
	}

	return warnings
}

// parameterValidator provides functionalities to check data against its parameter schema.
type parameterValidator struct {
	parameter *Parameter
}

func newParameterValidator(parameter *Parameter) *parameterValidator {
	return &parameterValidator{parameter}
}

// Validate validates value by comparing to its parameter schema.
func (v *parameterValidator) Validate(value interface{}) *ParameterWarning {
	if value == nil {
		if v.parameter.Optional {
			return nil
		}
		return v.newParameterWarning("required")
	}

	return v.validateType(value)
}

// validateType checks if value comforts its expected type.
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

// newParameterWarning creates a ParameterWarning with given warning message.
func (v *parameterValidator) newParameterWarning(warning string) *ParameterWarning {
	return &ParameterWarning{
		Key:       v.parameter.Key,
		Warning:   warning,
		Parameter: v.parameter,
	}
}
