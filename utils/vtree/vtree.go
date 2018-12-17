// Package vtree provides utility to analyze any Go type to print a type tree.
// type tree represents commonly used types by data-interchange formats.
// it resolves nested types like maps and structs until reach to the most basic
// types like numbers, strings and booleans.
package vtree

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// Type of a value.
type Type int

const (
	// Unknown type.
	Unknown Type = iota

	// Nil represents unset value.
	Nil

	// String type.
	String

	// Number type.
	Number

	// Bool type.
	Bool

	// Array type.
	// Corresponds to Go slices.
	Array

	// Object type.
	// Corresponds to Go maps and structs.
	Object
)

// Value represents a value.
type Value struct {
	// Type of the value.
	Type Type

	// Key is struct field name or map key.
	// It is filled when parent's type is object.
	Key string

	// Values are the child values inside the value.
	// It is filled when type is object or array.
	Values []Value
}

// GetByKey returns a child value of an object type by key.
func (v Value) GetByKey(key string, caseSensitive bool) (vv Value, ok bool) {
	for _, vv := range v.Values {
		vvKey := vv.Key
		key := key
		if !caseSensitive {
			vvKey = strings.ToLower(vvKey)
			key = strings.ToLower(key)
		}
		if vvKey == key {
			return vv, true
		}
	}
	return Value{}, false
}

// Analyze analyzes a Go type and produces a value tree.
func Analyze(v interface{}) Value {
	vv := Value{}
	if v == nil {
		return vv
	}
	return analyze(vv, reflect.ValueOf(v))
}

func analyze(v Value, rv reflect.Value) Value {
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface:
		return analyze(v, rv.Elem())

	case reflect.Map:
		v.Type = Object
		v.Values = make([]Value, 0)

		var (
			keys = []string{}
			// key-reflect.Value pair.
			values = map[string]reflect.Value{}
		)

		for _, key := range rv.MapKeys() {
			k := key.String()
			values[k] = key
			keys = append(keys, k)
		}
		// sort by map key to always keep the members of Values
		// slice in the same order.
		sort.Strings(keys)

		for _, key := range keys {
			v.Values = append(v.Values, analyze(Value{Key: key}, rv.MapIndex(values[key])))
		}

	case reflect.Struct:
		if _, ok := rv.Interface().(fmt.Stringer); ok {
			v.Type = String
		} else {
			v.Type = Object
			v.Values = make([]Value, 0)

			for i := 0; i < rv.NumField(); i++ {
				v.Values = append(v.Values, analyze(Value{Key: rv.Type().Field(i).Name}, rv.Field(i)))
			}
		}

	case reflect.Slice:
		v.Type = Array
		v.Values = make([]Value, 0)

		for i := 0; i < rv.Len(); i++ {
			v.Values = append(v.Values, analyze(Value{}, rv.Index(i)))
		}

	case reflect.String:
		v.Type = String

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		v.Type = Number

	case reflect.Bool:
		v.Type = Bool

	case reflect.Invalid:
		v.Type = Nil

	default:
		v.Type = Unknown
	}

	return v
}
