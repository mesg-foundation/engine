package ownership

import "github.com/mesg-foundation/engine/hash/hashserializer"

// HashSerialize returns the hashserialized string of this type
func (data *Ownership) HashSerialize() string {
	if data == nil {
		return ""
	}
	return hashserializer.New().
		AddString("2", data.Owner).
		AddString("3", data.ResourceHash.String()).
		HashSerialize()
}
