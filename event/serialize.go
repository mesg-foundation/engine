package event

import "github.com/mesg-foundation/engine/hash/hashserializer"

// HashSerialize returns the hashserialized string of this type
func (data *Event) HashSerialize() string {
	if data == nil {
		return ""
	}
	return hashserializer.New().
		AddString("2", data.InstanceHash.String()).
		AddString("3", data.Key).
		Add("4", data.Data).
		HashSerialize()
}
