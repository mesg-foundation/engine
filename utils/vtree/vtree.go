// Package vtree provides utility to analyse any Go types to print a type schema tree.
// it also resolves arrays, structs including nested types until reaching to basic
// types like numbers, strings and booleans.
// printed tree is meant to represent types that used by data encoding formats.
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
	Unknown Type = iota
	Nil
	String
	Number
	Bool
	Array
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

// GetByKey returns a child value by key of an object type.
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

// Analyse analyses a Go type and produces a value tree.
func Analyse(v interface{}) Value {
	if v == nil {
		return Value{}
	}
	return analyse(reflect.ValueOf(v))
}

func analyse(rv reflect.Value) Value {
	v := Value{}

	if rv.Kind() != reflect.Invalid {
		if _, ok := rv.Interface().(fmt.Stringer); ok {
			v.Type = String
			return v
		}
	}

	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface:
		return analyse(rv.Elem())

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
			val := analyse(rv.MapIndex(values[key]))
			val.Key = key
			v.Values = append(v.Values, val)
		}

	case reflect.Struct:
		tv := rv.Type()

		v.Type = Object
		v.Values = make([]Value, 0)

		for i := 0; i < rv.NumField(); i++ {
			val := analyse(rv.Field(i))
			val.Key = tv.Field(i).Name
			v.Values = append(v.Values, val)
		}

	case reflect.Slice:
		v.Type = Array
		v.Values = make([]Value, 0)

		for i := 0; i < rv.Len(); i++ {
			v.Values = append(v.Values, analyse(rv.Index(i)))
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
