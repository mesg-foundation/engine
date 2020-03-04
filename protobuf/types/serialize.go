package types

import (
	"sort"
	"strconv"

	"github.com/mesg-foundation/engine/hash/hashserializer"
)

// HashSerialize returns the hashserialized string of this type
func (data *Struct) HashSerialize() string {
	if data == nil {
		return ""
	}
	ser := hashserializer.New()
	ser.Add("1", mapValue(data.Fields))
	return ser.HashSerialize()
}

type mapValue map[string]*Value

// HashSerialize returns the hashserialized string of this type
func (data mapValue) HashSerialize() string {
	ser := hashserializer.New()
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		ser.Add(key, data[key])
	}
	return ser.HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Value) HashSerialize() string {
	if data == nil {
		return ""
	}
	ser := hashserializer.New()
	ser.AddFloat("2", data.GetNumberValue())
	ser.AddString("3", data.GetStringValue())
	ser.AddBool("4", data.GetBoolValue())
	ser.Add("5", data.GetStructValue())
	ser.Add("6", data.GetListValue())
	return ser.HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *ListValue) HashSerialize() string {
	if data == nil {
		return ""
	}
	ser := hashserializer.New()
	ser.Add("1", values(data.Values))
	return ser.HashSerialize()
}

type values []*Value

// HashSerialize returns the hashserialized string of this type
func (data values) HashSerialize() string {
	if data == nil {
		return ""
	}
	ser := hashserializer.New()
	for i, value := range data {
		ser.Add(strconv.Itoa(i), value)
	}
	return ser.HashSerialize()
}
