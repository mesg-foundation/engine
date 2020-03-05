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
	return hashserializer.New().
		Add("1", mapValue(data.Fields)).
		HashSerialize()
}

type mapValue map[string]*Value

// HashSerialize returns the hashserialized string of this type
func (data mapValue) HashSerialize() string {
	if data == nil || len(data) == 0 {
		return ""
	}
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	ser := hashserializer.New()
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
	return hashserializer.New().
		AddInt("1", int(data.GetNullValue())).
		AddFloat("2", data.GetNumberValue()).
		AddString("3", data.GetStringValue()).
		AddBool("4", data.GetBoolValue()).
		Add("5", data.GetStructValue()).
		Add("6", data.GetListValue()).
		HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *ListValue) HashSerialize() string {
	if data == nil {
		return ""
	}
	return hashserializer.New().
		Add("1", values(data.Values)).
		HashSerialize()
}

type values []*Value

// HashSerialize returns the hashserialized string of this type
func (data values) HashSerialize() string {
	if data == nil || len(data) == 0 {
		return ""
	}
	ser := hashserializer.New()
	for i, value := range data {
		ser.Add(strconv.Itoa(i), value)
	}
	return ser.HashSerialize()
}
