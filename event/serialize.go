package event

import "github.com/mesg-foundation/engine/hash/hashserializer"

// HashSerialize returns the hashserialized string of this type
func (data *Event) HashSerialize() string {
	if data == nil {
		return ""
	}
	ser := hashserializer.New()
	ser.AddString("2", data.InstanceHash.String())
	ser.AddString("3", data.Key)
	ser.Add("4", data.Data)
	return ser.HashSerialize()
}
