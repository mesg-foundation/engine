// Package hashserializer creates a string representing arbitrary value that can then be used to generate a deterministic hash.
//
// Each value must have a unique prefix.
//
// Empty or default value are discarded.
//
// The order to add the value should always be the same if you want the created string to be deterministic.
//
// Each value can be added to the HashSerializer using one of the AddXXX function alongside the prefix:
//  New().
//    AddBool(prefix, value).
//    AddFloat(prefix, value).
//    AddInt(prefix, value).
//    AddString(prefix, value).
//    AddStringSlice(prefix, value).
//    Add(prefix, value) // value must implement function `HashSerialize() string`.
//    HashSerialize()
//
// Format
//
// The format of the generated string for one value with a prefix is:
//  prefix:value;
// where value is the string representation of the value.
//
// When adding multiple values, they are append together:
//  prefix1:value1;prefix2:value2;
//
// Empty and default value
//
// Empty and default values are discarded and not part of the generated string. For slices and maps, if every value is empty then the slice or maps is also considered as empty.
//
// The following values are all considered as empty and will not be part of the generated string:
//  string("")
//  bool(false)
//  int(0)
//  float(0.0)
//  []string{}
//  []string{""}
//  map[string]string{}
//  map[string]string{"a":"", "b":""}
//
// String
//
// String values are the simplest are the value doesn't need to be encoded.
//
// For prefix "foo" and value "bar", the result is:
//  foo:bar;
//
// Bool
//
// Only bool value true produce an output as the default value false is discarded.
//
// For prefix "foo" and value "true", the result is:
//  foo:true;
//
// Int
//
// Int values are encoded to string using base 10 representation. If the value is negative, the sign "-" should prefix the value.
//
// For prefix "foo" and value "42", the result is:
//  foo:42;
//
// Float
//
// Float values are encoded to string using base 10 with no exponent. If the value is negative, the sign "-" should prefix the value.
//
// For prefix "foo" and value "3.14159265359", the result is:
//  foo:3.14159265359;
//
// Slice
//
// For slices, prefixes are the index of the elements. Empty value are discarded.
//
//  []string{"foo", "", "bar"}
//  0:foo;2:bar;
//
// Map
//
// For maps, prefixes are the keys. The maps must be sorted by keys to get a deterministic result. Empty value are discarded.
//
//  map[string]string{"c":"bar", "b": "", "a":"foo"}
//  a:foo;c:bar;
//
// Struct
//
// Struct should implement the function `HashSerialize() string` to return the hash-serialized string.
//
// Nested
//
// When adding a struct, a slice or a map, the value is hash-serialized before added to the string.
//
// For prefix "foo" and value `[]string{"foo", "", "bar"}`, the result is:
//  foo:0:foo;2:bar;;
//
//
package hashserializer

import (
	"bytes"
	"strconv"
)

// HashSerializable represents a serializable struct.
type HashSerializable interface {
	HashSerialize() string
}

// HashSerializer helps to create hash-serialized string of a struct.
type HashSerializer struct {
	buf bytes.Buffer
}

// New returns a new HashSerializer.
func New() *HashSerializer {
	return &HashSerializer{}
}

// HashSerialize returns the hash-serialized constructed string.
func (s *HashSerializer) HashSerialize() string {
	return s.buf.String()
}

// AddString adds a string to the serializer.
// If value is empty, nothing is added to the HashSerializer.
func (s *HashSerializer) AddString(prefix string, value string) *HashSerializer {
	if value != "" {
		s.buf.WriteString(prefix)
		s.buf.WriteString(":")
		s.buf.WriteString(value)
		s.buf.WriteString(";")
	}
	return s
}

// AddBool adds a boolean to the serializer.
// If value is false, nothing is added to the HashSerializer.
func (s *HashSerializer) AddBool(prefix string, value bool) *HashSerializer {
	if value {
		s.AddString(prefix, "true")
	}
	return s
}

// AddFloat adds a float to the serializer.
// If value is 0, nothing is added to the HashSerializer.
func (s *HashSerializer) AddFloat(prefix string, value float64) *HashSerializer {
	if value != 0 {
		s.AddString(prefix, strconv.FormatFloat(value, 'f', -1, 64))
	}
	return s
}

// AddInt adds a int to the serializer.
// If value is 0, nothing is added to the HashSerializer.
func (s *HashSerializer) AddInt(prefix string, value int) *HashSerializer {
	if value != 0 {
		s.AddString(prefix, strconv.Itoa(value))
	}
	return s
}

// Add adds a struct implementing HashSerializable to the serializer.
// The function HashSerialize will be called on value and its result injected as a string in HashSerializer.
// If value.HashSerialize() returns an empty string, nothing is added to the HashSerializer.
func (s *HashSerializer) Add(prefix string, value HashSerializable) *HashSerializer {
	s.AddString(prefix, value.HashSerialize())
	return s
}

// AddStringSlice adds a slice of string to the serializer.
// If value is an empty slice or each element is empty, nothing is added to the HashSerializer.
func (s *HashSerializer) AddStringSlice(prefix string, value []string) *HashSerializer {
	s.Add(prefix, StringSlice(value))
	return s
}

// StringSlice is an helper type to easily hash-serialized slice of string.
type StringSlice []string

// HashSerialize returns the hash-serialized string of this type.
func (s StringSlice) HashSerialize() string {
	ser := New()
	for i, value := range s {
		ser.AddString(strconv.Itoa(i), value)
	}
	return ser.HashSerialize()
}
