package service

import (
	"fmt"
	"math/big"
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
func validateParametersSchema(parameters []*Parameter, data map[string]interface{}) []*ParameterWarning {
	warningResults := make([]*ParameterWarning, 0, len(parameters))

	for _, param := range parameters {
		warnings := newParameterValidator(param).Validate(data[param.Key])
		warningResults = append(warningResults, warnings...)
	}

	return warningResults
}

// parameterValidator provides functionalities to check data against its parameter schema.
type parameterValidator struct {
	parameter *Parameter
}

func newParameterValidator(parameter *Parameter) *parameterValidator {
	return &parameterValidator{parameter}
}

// Validate validates value by comparing to its parameter schema.
func (v *parameterValidator) Validate(value interface{}) []*ParameterWarning {
	if value == nil {
		if v.parameter.Optional {
			return nil
		}
		return []*ParameterWarning{v.newParameterWarning("required")}
	}
	if v.parameter.Repeated {
		// Check if the value is a slice
		array, ok := value.([]interface{})
		if !ok {
			return []*ParameterWarning{v.newParameterWarning("not an array")}
		}
		for _, x := range array {
			if warnings := v.validateType(x); warnings != nil {
				return warnings
			}
		}
		return nil
	}
	return v.validateType(value)
}

// validateType checks if value comforts its expected type.
func (v *parameterValidator) validateType(value interface{}) []*ParameterWarning {
	switch v.parameter.Type {
	case "String":
		if _, ok := value.(string); !ok {
			return []*ParameterWarning{v.newParameterWarning("not a string")}
		}

	case "Number":
		switch value.(type) {
		case
			uint8, uint16, uint32, uint64, uint,
			int8, int16, int32, int64, int,
			float64, float32,
			*big.Int, *big.Float:
		default:
			return []*ParameterWarning{v.newParameterWarning("not a number")}
		}
	case "Boolean":
		if _, ok := value.(bool); !ok {
			return []*ParameterWarning{v.newParameterWarning("not a boolean")}
		}

	case "Object":
		data, okObj := value.(map[string]interface{})
		if !okObj {
			return []*ParameterWarning{v.newParameterWarning("not an object")}
		}
		return validateParametersSchema(v.parameter.Object, data)
	case "Any":
		return nil
	default:
		return []*ParameterWarning{v.newParameterWarning("an invalid type")}
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
