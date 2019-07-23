package convert

import (
	"reflect"

	"github.com/mesg-foundation/engine/protobuf/types"
)

// PbStructToMap converts protobuf struct to map[string]interface{}.
func PbStructToMap(s *types.Struct) map[string]interface{} {
	m := make(map[string]interface{})
	for k, v := range s.Fields {
		m[k] = PbValueToInterface(v)
	}
	return m
}

// PbValueToInterface converts protobuf value to interface{}.
func PbValueToInterface(v *types.Value) interface{} {
	switch v.Kind.(type) {
	case *types.Value_NullValue:
		return nil
	case *types.Value_NumberValue:
		return v.GetNumberValue()
	case *types.Value_StringValue:
		return v.GetStringValue()
	case *types.Value_BoolValue:
		return v.GetBoolValue()
	case *types.Value_StructValue:
		return PbStructToMap(v.GetStructValue())
	case *types.Value_ListValue:
		lv := v.GetListValue()
		if len(lv.Values) == 0 {
			return nil
		}
		a := make([]interface{}, len(lv.Values))
		for i, v := range lv.Values {
			a[i] = PbValueToInterface(v)
		}
		return a
	}
	return nil
}

// MapToPbStruct converts map[string]interface{} to protobuf struct.
func MapToPbStruct(m map[string]interface{}) *types.Struct {
	if len(m) == 0 {
		return nil
	}

	s := &types.Struct{
		Fields: make(map[string]*types.Value, len(m)),
	}

	for k, v := range m {
		s.Fields[k] = InterfaceToPbValue(v)
	}
	return s
}

// InterfaceToPbValue converts interface{} to protobuf value.
func InterfaceToPbValue(v interface{}) *types.Value {
	switch v := v.(type) {
	case nil:
		return nil
	case bool:
		return &types.Value{
			Kind: &types.Value_BoolValue{
				BoolValue: v,
			},
		}
	case int:
		return &types.Value{
			Kind: &types.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case int8:
		return &types.Value{
			Kind: &types.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case int16:
		return &types.Value{
			Kind: &types.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case int32:
		return &types.Value{
			Kind: &types.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case int64:
		return &types.Value{
			Kind: &types.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint:
		return &types.Value{
			Kind: &types.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint8:
		return &types.Value{
			Kind: &types.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint16:
		return &types.Value{
			Kind: &types.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint32:
		return &types.Value{
			Kind: &types.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint64:
		return &types.Value{
			Kind: &types.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case float32:
		return &types.Value{
			Kind: &types.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case float64:
		return &types.Value{
			Kind: &types.Value_NumberValue{
				NumberValue: v,
			},
		}
	case string:
		return &types.Value{
			Kind: &types.Value_StringValue{
				StringValue: v,
			},
		}
	case error:
		return &types.Value{
			Kind: &types.Value_StringValue{
				StringValue: v.Error(),
			},
		}
	default:
		return reflectValueToPbValue(reflect.ValueOf(v))
	}
}

// reflectValueToPbValue converts reflect.Value to protobuf value.
func reflectValueToPbValue(v reflect.Value) *types.Value {
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return nil
		}
		return reflectValueToPbValue(reflect.Indirect(v))
	case reflect.Array, reflect.Slice:
		if v.Len() == 0 {
			return nil
		}
		values := make([]*types.Value, v.Len())
		for i := 0; i < v.Len(); i++ {
			values[i] = reflectValueToPbValue(v.Index(i))
		}
		return &types.Value{
			Kind: &types.Value_ListValue{
				ListValue: &types.ListValue{
					Values: values,
				},
			},
		}
	case reflect.Struct:
		t := v.Type()
		size := v.NumField()
		if size == 0 {
			return nil
		}
		fields := make(map[string]*types.Value, size)
		for i := 0; i < size; i++ {
			name := t.Field(i).Name
			// Better way?
			if len(name) > 0 && 'A' <= name[0] && name[0] <= 'Z' {
				fields[name] = reflectValueToPbValue(v.Field(i))
			}
		}
		if len(fields) == 0 {
			return nil
		}
		return &types.Value{
			Kind: &types.Value_StructValue{
				StructValue: &types.Struct{
					Fields: fields,
				},
			},
		}
	case reflect.Map:
		keys := v.MapKeys()
		if len(keys) == 0 {
			return nil
		}
		fields := make(map[string]*types.Value, len(keys))
		for _, k := range keys {
			if k.Kind() == reflect.String {
				fields[k.String()] = reflectValueToPbValue(v.MapIndex(k))
			}
		}
		if len(fields) == 0 {
			return nil
		}
		return &types.Value{
			Kind: &types.Value_StructValue{
				StructValue: &types.Struct{
					Fields: fields,
				},
			},
		}
	default:
		return nil
	}
}
