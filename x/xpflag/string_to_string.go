// Copyright 2018 MESG Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package xpflag

import (
	"bytes"
	"fmt"
	"strings"
)

// StringToStringValue flag.
type StringToStringValue struct {
	value     *map[string]string
	separator string
	changed   bool
}

// NewStringToStringValue creates new flag with init map, default value and comma separator.
func NewStringToStringValue(p *map[string]string, value map[string]string) *StringToStringValue {
	s := new(StringToStringValue)
	s.separator = ","
	s.value = p
	*s.value = value
	return s
}

// SetSeparator changes separator.
func (s *StringToStringValue) SetSeparator(separator string) {
	s.separator = separator
}

// Set value in format: a=1,b=2
func (s *StringToStringValue) Set(val string) error {
	ss := strings.Split(val, s.separator)
	out := make(map[string]string, len(ss))
	for _, pair := range ss {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			return fmt.Errorf("%s must be formatted as key=value", pair)
		}
		out[kv[0]] = kv[1]
	}
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
			buf.WriteString(s.separator)
		}
		buf.WriteString(k)
		buf.WriteRune('=')
		buf.WriteString(v)
		i++
	}
	return buf.String()
}
