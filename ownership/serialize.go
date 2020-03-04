package ownership

import "github.com/mesg-foundation/engine/hash/hashserializer"

// HashSerialize returns the hashserialized string of this type
func (data *Ownership) HashSerialize() string {
	if data == nil {
		return ""
	}
	ser := hashserializer.New()
	ser.AddString("2", data.Owner)
	ser.AddString("3", data.ResourceHash.String())
	return ser.HashSerialize()
}
