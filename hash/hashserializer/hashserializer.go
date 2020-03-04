package hashserializer

import (
	"bytes"
	"strconv"
)

type HashSerializable interface {
	HashSerialize() string
}

type HashSerializer struct {
	buf bytes.Buffer
}

// --------------------------------------------------------

func New() *HashSerializer {
	return &HashSerializer{}
}

func (s *HashSerializer) HashSerialize() string {
	return s.buf.String()
}

func (s *HashSerializer) AddString(number string, value string) {
	if value != "" {
		s.buf.WriteString(number)
		s.buf.WriteString(":")
		s.buf.WriteString(value)
		s.buf.WriteString(";")
	}
}

func (s *HashSerializer) AddBool(number string, value bool) {
	if value {
		s.AddString(number, "true")
	}
}

func (s *HashSerializer) AddFloat(number string, value float64) {
	if value != 0 {
		s.AddString(number, strconv.FormatFloat(value, 'f', -1, 64))
	}
}

func (s *HashSerializer) AddInt(number string, value int) {
	if value != 0 {
		s.AddString(number, strconv.Itoa(value))
	}
}

func (s *HashSerializer) Add(number string, value HashSerializable) {
	s.AddString(number, value.HashSerialize())
}

func (s *HashSerializer) AddStringSlice(number string, value []string) {
	s.Add(number, StringSlice(value))
}

// --------------------------------------------------------

type StringSlice []string

func (s StringSlice) HashSerialize() string {
	ser := New()
	for i, value := range s {
		ser.AddString(strconv.Itoa(i), value)
	}
	return ser.HashSerialize()
}
