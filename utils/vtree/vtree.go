// Package vtree provides utility to analyze Go types to print a type tree.
// type tree represents commonly used types by data-interchange formats.
// it's capable to resolve nested maps until reach to the most basic
// types like numbers, strings and booleans.
package vtree

import (
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
	// Corresponds to Go maps.
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
// v can be;
// * basic types like string, number types, bool;
// * slices of string, number types, bool, maps for one level (e.g. []string and not [][]string);
// * maps & nested maps with map[string]interface{} type which can contain basic types and slices.
// pointer values and structs aren't supported.
func Analyze(v interface{}) Value {
	return analyze(Value{}, v)
}

func analyze(vv Value, v interface{}) Value {
	switch d := v.(type) {
	case map[string]interface{}:
		vv.Type = Object
		vv.Values = make([]Value, 0)

		keys := []string{}
		for key := range d {
			keys = append(keys, key)
		}
		// sort by map keys to always keep the members of Values
		// slice in the same order.
		sort.Strings(keys)

		for _, key := range keys {
			vv.Values = append(vv.Values, analyze(Value{Key: key}, d[key]))
		}

	case []map[string]interface{}:
		vv.Type = Array
		vv.Values = make([]Value, 0)

		for _, val := range d {
			vv.Values = append(vv.Values, analyze(Value{}, val))
		}

	case []interface{}:
		vv.Type = Array
		vv.Values = make([]Value, 0)

		for _, val := range d {
			vv.Values = append(vv.Values, analyze(Value{}, val))
		}

	case []string:
		vv.Type = Array
		vv.Values = createValues(String, len(d))

	case []int:
		vv.Type = Array
		vv.Values = createValues(Number, len(d))

	case []int8:
		vv.Type = Array
		vv.Values = createValues(Number, len(d))

	case []int16:
		vv.Type = Array
		vv.Values = createValues(Number, len(d))

	case []int32:
		vv.Type = Array
		vv.Values = createValues(Number, len(d))

	case []int64:
		vv.Type = Array
		vv.Values = createValues(Number, len(d))

	case []uint:
		vv.Type = Array
		vv.Values = createValues(Number, len(d))

	case []uint8:
		vv.Type = Array
		vv.Values = createValues(Number, len(d))

	case []uint16:
		vv.Type = Array
		vv.Values = createValues(Number, len(d))

	case []uint32:
		vv.Type = Array
		vv.Values = createValues(Number, len(d))

	case []uint64:
		vv.Type = Array
		vv.Values = createValues(Number, len(d))

	case []float32:
		vv.Type = Array
		vv.Values = createValues(Number, len(d))

	case []float64:
		vv.Type = Array
		vv.Values = createValues(Number, len(d))

	case []bool:
		vv.Type = Array
		vv.Values = createValues(Bool, len(d))

	case string:
		vv.Type = String

	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		vv.Type = Number

	case bool:
		vv.Type = Bool

	case nil:
		vv.Type = Nil

	default:
		vv.Type = Unknown
	}

	return vv
}

// createValues creates a slice of n values with the same t.
func createValues(t Type, n int) []Value {
	values := make([]Value, 0)
	for i := 0; i < n; i++ {
		values = append(values, Value{Type: t})
	}
	return values
}
