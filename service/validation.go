package service

import (
	"errors"
	"fmt"

	"github.com/mesg-foundation/core/utils/vtree"
)

type parameterValidator struct {
	warnings []*ParameterWarning
}

func newParameterValidator() *parameterValidator {
	return &parameterValidator{warnings: make([]*ParameterWarning, 0)}
}

// Validate validates data to check if it matches with parameter's schema.
// TODO(ilgooz) add this as a method to Service type when create custom types for Event, Task etc.
// TODO(ilgooz) remove pointer from *Parameter.
func (p *parameterValidator) Validate(params []*Parameter, data interface{}) ([]*ParameterWarning, error) {
	v := vtree.Analyse(data)
	if v.Type != vtree.Object {
		return nil, errors.New("not an object")
	}
	p.validateParameters(params, v)
	return p.warnings, nil
}

func (p *parameterValidator) validateParameters(params []*Parameter, v vtree.Value) {
	for _, param := range params {
		v, ok := v.GetByKey(param.Key, false)
		if !param.Optional && !ok {
			p.addWarning(param, "required")
			continue
		}
		p.validateParameter(param, v)
	}
}

func (p *parameterValidator) validateParameter(param *Parameter, v vtree.Value) {
	if param.Repeated {
		if v.Type == vtree.Array {
			for _, vv := range v.Values {
				p.validateParameterType(param, vv)
			}
		} else {
			p.addWarning(param, "not an array")
		}
	} else {
		p.validateParameterType(param, v)
	}
}

func (p *parameterValidator) validateParameterType(param *Parameter, v vtree.Value) {
	if !param.Optional && v.Type == vtree.Nil {
		p.addWarning(param, "required")
		return
	}

	switch {
	// interpret basic types.
	case param.Type != "":
		switch param.Type {
		case "String":
			if v.Type != vtree.String {
				p.addWarning(param, "not a string")
			}
		case "Number":
			if v.Type != vtree.Number {
				p.addWarning(param, "not a number")
			}
		case "Boolean":
			if v.Type != vtree.Bool {
				p.addWarning(param, "not a boolean")
			}
		// TODO(ilgooz) Deprecate Object in future.
		case "Any", "Object":
		default:
			p.addWarning(param, "unknown type")
		}

	// interpret nested parameters.
	default:
		if v.Type != vtree.Object {
			p.addWarning(param, "not an object")
		} else {
			p.validateParameters(param.Parameters, v)
		}
	}
}

func (p *parameterValidator) addWarning(param *Parameter, warning string) {
	p.warnings = append(p.warnings, &ParameterWarning{
		Key:       param.Key,
		Warning:   warning,
		Parameter: param,
	})
}

// ParameterWarning contains a specific warning related to a parameter.
type ParameterWarning struct {
	Key       string
	Warning   string
	Parameter *Parameter
}

func (p *ParameterWarning) String() string {
	return fmt.Sprintf("Value of %q is %s", p.Key, p.Warning)
}
