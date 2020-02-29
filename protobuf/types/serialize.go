package types

import (
	"sort"
	"strconv"

	"github.com/mesg-foundation/engine/hash/serializer"
)

func (data *Struct) Serialize() string {
	if data == nil {
		return ""
	}
	ser := serializer.New()
	ser.Add("1", mapValue(data.Fields))
	return ser.Serialize()
}

type mapValue map[string]*Value

func (data mapValue) Serialize() string {
	ser := serializer.New()
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		ser.Add(key, data[key])
	}
	return ser.Serialize()
}

func (data *Value) Serialize() string {
	if data == nil {
		return ""
	}
	ser := serializer.New()
	ser.AddFloat("2", data.GetNumberValue())
	ser.AddString("3", data.GetStringValue())
	ser.AddBool("4", data.GetBoolValue())
	ser.Add("5", data.GetStructValue())
	ser.Add("6", data.GetListValue())
	return ser.Serialize()
}

func (data *ListValue) Serialize() string {
	if data == nil {
		return ""
	}
	ser := serializer.New()
	ser.Add("1", values(data.Values))
	return ser.Serialize()
}

type values []*Value

func (data values) Serialize() string {
	if data == nil {
		return ""
	}
	ser := serializer.New()
	for i, value := range data {
		ser.Add(strconv.Itoa(i), value)
	}
	return ser.Serialize()
}
