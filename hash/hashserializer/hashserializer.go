// Package hashserializer creates a string representing arbitrary value that can then be used to generate a deterministic hash.
//
// Each value can be added to the HashSerializer using one of the AddXXX function alongside a prefix.
// The prefix should be unique and the AddXXX function should always be called in the same order in order to create a deterministic string.
//
// Empty value are discarded and its prefix will not be part of the generated string.
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
func (s *HashSerializer) AddString(prefix string, value string) {
	if value != "" {
		s.buf.WriteString(prefix)
		s.buf.WriteString(":")
		s.buf.WriteString(value)
		s.buf.WriteString(";")
	}
}

// AddBool adds a boolean to the serializer.
// If value is false, nothing is added to the HashSerializer.
func (s *HashSerializer) AddBool(prefix string, value bool) {
	if value {
		s.AddString(prefix, "true")
	}
}

// AddFloat adds a float to the serializer.
// If value is 0, nothing is added to the HashSerializer.
func (s *HashSerializer) AddFloat(prefix string, value float64) {
	if value != 0 {
		s.AddString(prefix, strconv.FormatFloat(value, 'f', -1, 64))
	}
}

// AddInt adds a int to the serializer.
// If value is 0, nothing is added to the HashSerializer.
func (s *HashSerializer) AddInt(prefix string, value int) {
	if value != 0 {
		s.AddString(prefix, strconv.Itoa(value))
	}
}

// Add adds a struct implementing HashSerializable to the serializer.
// The function HashSerialize will be called on value and its result injected as a string in HashSerializer.
// If value.HashSerialize() returns an empty string, nothing is added to the HashSerializer.
func (s *HashSerializer) Add(prefix string, value HashSerializable) {
	s.AddString(prefix, value.HashSerialize())
}

// AddStringSlice adds a slice of string to the serializer.
// If value is an empty slice or each element is empty, nothing is added to the HashSerializer.
func (s *HashSerializer) AddStringSlice(prefix string, value []string) {
	s.Add(prefix, StringSlice(value))
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
