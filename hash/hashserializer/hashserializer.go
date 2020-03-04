package hashserializer

import (
	"bytes"
	"strconv"
)

// HashSerializable represents a serializable struct
type HashSerializable interface {
	HashSerialize() string
}

// HashSerializer helps to create hash-serialized string of a struct
type HashSerializer struct {
	buf bytes.Buffer
}

// New returns a new HashSerializer
func New() *HashSerializer {
	return &HashSerializer{}
}

// HashSerialize returns the hash-serialized constructed string
func (s *HashSerializer) HashSerialize() string {
	return s.buf.String()
}

// AddString adds a string to the serializer
func (s *HashSerializer) AddString(number string, value string) {
	if value != "" {
		s.buf.WriteString(number)
		s.buf.WriteString(":")
		s.buf.WriteString(value)
		s.buf.WriteString(";")
	}
}

// AddBool adds a boolean to the serializer
func (s *HashSerializer) AddBool(number string, value bool) {
	if value {
		s.AddString(number, "true")
	}
}

// AddFloat adds a float to the serializer
func (s *HashSerializer) AddFloat(number string, value float64) {
	if value != 0 {
		s.AddString(number, strconv.FormatFloat(value, 'f', -1, 64))
	}
}

// AddInt adds a int to the serializer
func (s *HashSerializer) AddInt(number string, value int) {
	if value != 0 {
		s.AddString(number, strconv.Itoa(value))
	}
}

// Add adds a hashserializable struct to the serializer
func (s *HashSerializer) Add(number string, value HashSerializable) {
	s.AddString(number, value.HashSerialize())
}

// AddStringSlice adds a slice of string to the serializer
func (s *HashSerializer) AddStringSlice(number string, value []string) {
	s.Add(number, StringSlice(value))
}

// StringSlice is an helper type to easily hashserialized slice of string
type StringSlice []string

// HashSerialize returns the hashserialized string of this type
func (s StringSlice) HashSerialize() string {
	ser := New()
	for i, value := range s {
		ser.AddString(strconv.Itoa(i), value)
	}
	return ser.HashSerialize()
}
