package service

// ParameterWarning is a struct that contains a specific warning related to a parameter
type ParameterWarning struct {
	warning   string
	parameter *Parameter
}

// TODO: Rename this file into validation

func validParameters(parameters map[string]*Parameter, data map[string]interface{}) bool {
	return len(validateParameters(parameters, data)) == 0
}

func validateParameters(parameters map[string]*Parameter, data map[string]interface{}) []*ParameterWarning {
	warnings := make([]*ParameterWarning, 0)
	for key, param := range parameters {
		warning := param.Validate(data[key])
		if warning != nil {
			warnings = append(warnings, warning)
		}
	}
	return warnings
}

func (p *ParameterWarning) String() string {
	return p.warning
}

// IsValid returns true if the data are valid for a specific event
func (e *Event) IsValid(data map[string]interface{}) bool {
	return validParameters(e.Data, data)
}

// Validate data for a specific event
func (e *Event) Validate(data map[string]interface{}) []*ParameterWarning {
	return validateParameters(e.Data, data)
}

// IsValid returns true if the inputs are valid for a specific task
func (t *Task) IsValid(data map[string]interface{}) bool {
	return validParameters(t.Inputs, data)
}

// Validate inputs for a specific task
func (t *Task) Validate(data map[string]interface{}) []*ParameterWarning {
	return validateParameters(t.Inputs, data)
}

// Validate returns a warning based on the match of the data given in parameter and the parameter
func (p *Parameter) Validate(data interface{}) *ParameterWarning {
	if data == nil {
		if p.Optional {
			return nil
		}
		return &ParameterWarning{parameter: p, warning: "required"}
	}
	value := data
	var warning *ParameterWarning
	switch p.Type {
	case "String":
		_, ok := value.(string)
		if !ok {
			warning = &ParameterWarning{parameter: p, warning: "not a string"}
		}
	case "Number":
		_, okFloat64 := value.(float64)
		_, okFloat32 := value.(float32)
		_, okInt := value.(int)
		if !okInt && !okFloat64 && !okFloat32 {
			warning = &ParameterWarning{parameter: p, warning: "not a number"}
		}
	case "Boolean":
		_, ok := value.(bool)
		if !ok {
			warning = &ParameterWarning{parameter: p, warning: "not a boolean"}
		}
	case "Object":
		_, okObj := value.(map[string]interface{})
		_, okArr := value.([]interface{})
		if !okObj && !okArr {
			warning = &ParameterWarning{parameter: p, warning: "not an object/array"}
		}
	default:
		warning = &ParameterWarning{parameter: p, warning: "invalid type"}
	}
	return warning
}
