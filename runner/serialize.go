package runner

import "github.com/mesg-foundation/engine/hash/hashserializer"

// HashSerialize returns the hashserialized string of this type
func (data *Runner) HashSerialize() string {
	if data == nil {
		return ""
	}
	return hashserializer.New().
		AddString("2", data.Address).
		AddString("3", data.InstanceHash.String()).
		HashSerialize()
}
