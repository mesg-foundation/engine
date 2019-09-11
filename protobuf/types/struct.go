package types

// NewStruct creates a new struct with given fields.
func NewStruct(fields map[string]*Value) *Struct {
	return &Struct{Fields: fields}
}

// GetNullValue returns the string value of key.
func (s *Struct) GetNullValue(key string) NullValue {
	if s.Fields == nil || s.Fields[key] == nil {
		return NullValue_NULL_VALUE
	}
	return s.Fields[key].GetNullValue()
}

// GetNumberValue returns the string value of key.
func (s *Struct) GetNumberValue(key string) float64 {
	if s.Fields == nil || s.Fields[key] == nil {
		return 0
	}
	return s.Fields[key].GetNumberValue()
}

// GetStringValue returns the string value of key.
func (s *Struct) GetStringValue(key string) string {
	if s.Fields == nil || s.Fields[key] == nil {
		return ""
	}
	return s.Fields[key].GetStringValue()
}

// GetBoolValue returns the string value of key.
func (s *Struct) GetBoolValue(key string) bool {
	if s.Fields == nil || s.Fields[key] == nil {
		return false
	}
	return s.Fields[key].GetBoolValue()
}

// GetStructValue returns the string value of key.
func (s *Struct) GetStructValue(key string) *Struct {
	if s.Fields == nil || s.Fields[key] == nil {
		return nil
	}
	return s.Fields[key].GetStructValue()
}

// GetListValue returns the string value of key.
func (s *Struct) GetListValue(key string) *ListValue {
	if s.Fields == nil || s.Fields[key] == nil {
		return nil
	}
	return s.Fields[key].GetListValue()
}

// NewValueFrom creates new value from underlying data.
func NewValueFrom(val interface{}) *Value {
	switch val := val.(type) {
	case nil:
		return &Value{
			Kind: &Value_NullValue{},
		}
	case *Value:
		return val
	case []*Value:
		return &Value{
			Kind: &Value_ListValue{
				ListValue: &ListValue{
					Values: val,
				},
			},
		}
	case *Value_NullValue, *Value_NumberValue,
		*Value_StringValue, *Value_BoolValue,
		*Value_StructValue, *Value_ListValue:
		return &Value{Kind: val.(isValue_Kind)}
	case bool:
		return &Value{
			Kind: &Value_BoolValue{val},
		}
	case *bool:
		return &Value{
			Kind: &Value_BoolValue{*val},
		}
	case int8:
		return &Value{
			Kind: &Value_NumberValue{float64(val)},
		}
	case uint8:
		return &Value{
			Kind: &Value_NumberValue{float64(val)},
		}
	case int16:
		return &Value{
			Kind: &Value_NumberValue{float64(val)},
		}
	case uint16:
		return &Value{
			Kind: &Value_NumberValue{float64(val)},
		}
	case int32:
		return &Value{
			Kind: &Value_NumberValue{float64(val)},
		}
	case uint32:
		return &Value{
			Kind: &Value_NumberValue{float64(val)},
		}
	case int64:
		return &Value{
			Kind: &Value_NumberValue{float64(val)},
		}
	case uint64:
		return &Value{
			Kind: &Value_NumberValue{float64(val)},
		}
	case float32:
		return &Value{
			Kind: &Value_NumberValue{float64(val)},
		}
	case float64:
		return &Value{
			Kind: &Value_NumberValue{val},
		}
	case *int8:
		return &Value{
			Kind: &Value_NumberValue{float64(*val)},
		}
	case *uint8:
		return &Value{
			Kind: &Value_NumberValue{float64(*val)},
		}
	case *int16:
		return &Value{
			Kind: &Value_NumberValue{float64(*val)},
		}
	case *uint16:
		return &Value{
			Kind: &Value_NumberValue{float64(*val)},
		}
	case *int32:
		return &Value{
			Kind: &Value_NumberValue{float64(*val)},
		}
	case *uint32:
		return &Value{
			Kind: &Value_NumberValue{float64(*val)},
		}
	case *int64:
		return &Value{
			Kind: &Value_NumberValue{float64(*val)},
		}
	case *uint64:
		return &Value{
			Kind: &Value_NumberValue{float64(*val)},
		}
	case *float32:
		return &Value{
			Kind: &Value_NumberValue{float64(*val)},
		}
	case *float64:
		return &Value{
			Kind: &Value_NumberValue{*val},
		}
	case []byte:
		return &Value{
			Kind: &Value_StringValue{string(val)},
		}
	case string:
		return &Value{
			Kind: &Value_StringValue{val},
		}
	case map[string]interface{}:
		s := &Struct{
			Fields: make(map[string]*Value),
		}
		for k, v := range val {
			s.Fields[k] = NewValueFrom(v)
		}
		return &Value{Kind: &Value_StructValue{s}}
	}
	panic("not supported")
}
