package service

import (
	"fmt"

	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/x/xerrors"
)

// validateServiceParameters validates data to see if it matches with parameters schema.
func validateServiceParameters(parameters []*Service_Parameter, data *types.Struct) error {
	var errs xerrors.Errors

	for _, p := range parameters {
		var value *types.Value
		if data != nil && data.Fields != nil {
			value = data.Fields[p.Key]
		}
		if err := p.Validate(value); err != nil {
			errs = append(errs, err)
		}
	}

	return errs.ErrorOrNil()
}

// Validate checks if service parameter hash proper types for arrays and objects.
func (p *Service_Parameter) Validate(value *types.Value) error {
	if value == nil {
		if p.Optional {
			return nil
		}
		return fmt.Errorf("value of %q is required", p.Key)
	}

	if p.Repeated {
		array := value.GetListValue()
		if array == nil {
			return fmt.Errorf("value of %q is not an array", p.Key)
		}

		for _, value := range array.Values {
			if err := p.validateType(value); err != nil {
				return err
			}
		}
		return nil
	}
	return p.validateType(value)
}

// validateType checks if value comforts its expected type.
func (p *Service_Parameter) validateType(value *types.Value) error {
	switch p.Type {
	case "String":
		if _, ok := value.GetKind().(*types.Value_StringValue); !ok {
			return fmt.Errorf("value of %q is not a string", p.Key)
		}
	case "Number":
		if _, ok := value.GetKind().(*types.Value_NumberValue); !ok {
			return fmt.Errorf("value of %q is not a number", p.Key)
		}
	case "Boolean":
		if _, ok := value.GetKind().(*types.Value_BoolValue); !ok {
			return fmt.Errorf("value of %q is not a boolean", p.Key)
		}
	case "Object":
		obj, ok := value.GetKind().(*types.Value_StructValue)
		if !ok {
			return fmt.Errorf("value of %q is not an object", p.Key)
		}
		return validateServiceParameters(p.Object, obj.StructValue)
	case "Any":
		return nil
	default:
		return fmt.Errorf("value of %q has an invalid type ", p.Key)
	}

	return nil
}
