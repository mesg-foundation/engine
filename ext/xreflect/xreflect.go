package xreflect

import (
	"reflect"
)

// IsNil reports whether its argument o is nil.
// It traverses down to first non-pointer value.
func IsNil(o interface{}) bool {
	if o == nil {
		return true
	}

	for rv := reflect.ValueOf(o); rv.Kind() == reflect.Ptr; rv = rv.Elem() {
		if rv.IsNil() {
			return true
		}
	}
	return false
}
