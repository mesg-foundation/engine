package convert

import (
	"fmt"
	"reflect"

	structpb "github.com/golang/protobuf/ptypes/struct"
)

// PbStructToMap converts protobuf struct to map[string]interface{}.
func PbStructToMap(s *structpb.Struct) map[string]interface{} {
	m := make(map[string]interface{})
	for k, v := range s.Fields {
		m[k] = PbValueToInterface(v)
	}
	return m
}

// PbValueToInterface converts protobuf value to interface{}.
func PbValueToInterface(v *structpb.Value) interface{} {
	if v == nil {
		return nil
	}
	switch v.Kind.(type) {
	case *structpb.Value_NullValue:
		return nil
	case *structpb.Value_NumberValue:
		return v.GetNumberValue()
	case *structpb.Value_StringValue:
		return v.GetStringValue()
	case *structpb.Value_BoolValue:
		return v.GetBoolValue()
	case *structpb.Value_StructValue:
		return PbStructToMap(v.GetStructValue())
	case *structpb.Value_ListValue:
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
func MapToPbStruct(m map[string]interface{}) *structpb.Struct {
	if len(m) == 0 {
		return nil
	}

	s := &structpb.Struct{
		Fields: make(map[string]*structpb.Value, len(m)),
	}

	for k, v := range m {
		s.Fields[k] = InterfaceToPbValue(v)
	}
	return s
}

// InterfaceToPbValue converts interface{} to protobuf value.
func InterfaceToPbValue(v interface{}) *structpb.Value {
	switch v := v.(type) {
	case nil:
		return nil
	case bool:
		return &structpb.Value{
			Kind: &structpb.Value_BoolValue{
				BoolValue: v,
			},
		}
	case int:
		return &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case int8:
		return &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case int16:
		return &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case int32:
		return &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case int64:
		return &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint:
		return &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint8:
		return &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint16:
		return &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint32:
		return &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case uint64:
		return &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case float32:
		return &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: float64(v),
			},
		}
	case float64:
		return &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: v,
			},
		}
	case string:
		return &structpb.Value{
			Kind: &structpb.Value_StringValue{
				StringValue: v,
			},
		}
	case error:
		return &structpb.Value{
			Kind: &structpb.Value_StringValue{
				StringValue: v.Error(),
			},
		}
	default:
		return reflectValueToPbValue(reflect.ValueOf(v))
	}
}

// reflectValueToPbValue converts reflect.Value to protobuf value.
func reflectValueToPbValue(v reflect.Value) *structpb.Value {
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return nil
		}
		return reflectValueToPbValue(reflect.Indirect(v))
	case reflect.Interface:
		return InterfaceToPbValue(v.Interface())
	case reflect.Bool:
		return &structpb.Value{
			Kind: &structpb.Value_BoolValue{
				BoolValue: v.Bool(),
			},
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: float64(v.Int()),
			},
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: float64(v.Uint()),
			},
		}
	case reflect.Float32, reflect.Float64:
		return &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: v.Float(),
			},
		}
	case reflect.String:
		return &structpb.Value{
			Kind: &structpb.Value_StringValue{
				StringValue: v.String(),
			},
		}
	case reflect.Array, reflect.Slice:
		if v.Len() == 0 {
			return nil
		}
		values := make([]*structpb.Value, v.Len())
		for i := 0; i < v.Len(); i++ {
			values[i] = reflectValueToPbValue(v.Index(i))
		}
		return &structpb.Value{
			Kind: &structpb.Value_ListValue{
				ListValue: &structpb.ListValue{
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
		fields := make(map[string]*structpb.Value, size)
		for i := 0; i < size; i++ {
			name := t.Field(i).Name
			if len(name) > 0 && 'A' <= name[0] && name[0] <= 'Z' {
				fields[name] = reflectValueToPbValue(v.Field(i))
			}
		}
		if len(fields) == 0 {
			return nil
		}
		return &structpb.Value{
			Kind: &structpb.Value_StructValue{
				StructValue: &structpb.Struct{
					Fields: fields,
				},
			},
		}
	case reflect.Map:
		keys := v.MapKeys()
		if len(keys) == 0 {
			return nil
		}
		fields := make(map[string]*structpb.Value, len(keys))
		for _, k := range keys {
			if k.Kind() == reflect.String {
				fields[k.String()] = reflectValueToPbValue(v.MapIndex(k))
			}
		}
		if len(fields) == 0 {
			return nil
		}
		return &structpb.Value{
			Kind: &structpb.Value_StructValue{
				StructValue: &structpb.Struct{
					Fields: fields,
				},
			},
		}
	default:
		return &structpb.Value{
			Kind: &structpb.Value_StringValue{
				StringValue: fmt.Sprint(v),
			},
		}
	}
}
