package serializer

import (
	"bytes"
	"strconv"
)

type Serializable interface {
	Serialize() string
}

type Serializer struct {
	buf bytes.Buffer
}

// --------------------------------------------------------

func New() *Serializer {
	return &Serializer{}
}

func (s *Serializer) Serialize() string {
	return s.buf.String()
}

func (s *Serializer) AddString(number string, value string) {
	if value != "" {
		s.buf.WriteString(number)
		s.buf.WriteString(":")
		s.buf.WriteString(value)
		s.buf.WriteString(";")
	}
}

func (s *Serializer) AddBool(number string, value bool) {
	if value {
		s.AddString(number, "true")
	}
}

func (s *Serializer) AddFloat(number string, value float64) {
	if value != 0 {
		s.AddString(number, strconv.FormatFloat(value, 'f', -1, 64))
	}
}

func (s *Serializer) AddInt(number string, value int) {
	if value != 0 {
		s.AddString(number, strconv.Itoa(value))
	}
}

func (s *Serializer) Add(number string, value Serializable) {
	s.AddString(number, value.Serialize())
}

func (s *Serializer) AddStringSlice(number string, value []string) {
	s.Add(number, StringSlice(value))
}

// --------------------------------------------------------

type StringSlice []string

func (s StringSlice) Serialize() string {
	ser := New()
	for i, value := range s {
		ser.AddString(strconv.Itoa(i), value)
	}
	return ser.Serialize()
}
