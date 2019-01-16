package xpflag

import (
	"bytes"
	"fmt"
	"strings"
)

// StringToStringValue flag.
type StringToStringValue struct {
	value   *map[string]string
	changed bool
}

// NewStringToStringValue creates new flag with init map, default value.
func NewStringToStringValue(p *map[string]string, value map[string]string) *StringToStringValue {
	s := new(StringToStringValue)
	s.value = p
	*s.value = value
	return s
}

// Set value in format: a=1,b=2
func (s *StringToStringValue) Set(val string) error {
	out := make(map[string]string)
	kv := strings.SplitN(val, "=", 2)
	if len(kv) != 2 {
		return fmt.Errorf("%s must be formatted as key=value", val)
	}
	out[kv[0]] = kv[1]
	if !s.changed {
		*s.value = out
	} else {
		for k, v := range out {
			(*s.value)[k] = v
		}
	}
	s.changed = true
	return nil
}

// Type returns type of value.
func (s *StringToStringValue) Type() string {
	return "key=value"
}

func (s *StringToStringValue) String() string {
	var buf bytes.Buffer
	i := 0
	for k, v := range *s.value {
		if i > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(k)
		buf.WriteRune('=')
		buf.WriteString(v)
		i++
	}
	return buf.String()
}
