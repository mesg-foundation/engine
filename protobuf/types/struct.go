package types

import (
	"sort"
)

// KeyValue is a simple key/value representation of one field of a Struct.
type KeyValue struct {
	Key   string
	Value *Value
}

// MarshalAmino transforms the Struct to an array of key/value.
func (m Struct) MarshalAmino() ([]KeyValue, error) {
	p := make([]KeyValue, len(m.Fields))
	fieldKeys := make([]string, len(m.Fields))
	i := 0
	for key := range m.Fields {
		fieldKeys[i] = key
		i++
	}
	sort.Stable(sort.StringSlice(fieldKeys))
	for i, key := range fieldKeys {
		p[i] = KeyValue{
			Key:   key,
			Value: m.Fields[key],
		}
	}
	return p, nil
}

// UnmarshalAmino transforms the key/value array to a Struct.
func (m *Struct) UnmarshalAmino(keyValues []KeyValue) error {
	m.Fields = make(map[string]*Value, len(keyValues))
	for _, p := range keyValues {
		m.Fields[p.Key] = p.Value
	}
	return nil
}
