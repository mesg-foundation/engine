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
	p := make([]KeyValue, 0)
	fieldKeys := make([]string, 0)
	for key := range m.Fields {
		fieldKeys = append(fieldKeys, key)
	}
	sort.Stable(sort.StringSlice(fieldKeys))
	for _, key := range fieldKeys {
		p = append(p, KeyValue{
			Key:   key,
			Value: m.Fields[key],
		})
	}
	return p, nil
}

// UnmarshalAmino transforms the key/value array to a Struct.
func (m *Struct) UnmarshalAmino(keyValues []KeyValue) error {
	fields := make(map[string]*Value)
	for _, p := range keyValues {
		fields[p.Key] = p.Value
	}
	m.Fields = fields
	return nil
}
