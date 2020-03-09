package execution

import "github.com/mesg-foundation/engine/hash/hashserializer"

// HashSerialize returns the hashserialized string of this type
func (data *Execution) HashSerialize() string {
	return hashserializer.New().
		AddString("2", data.ParentHash.String()).
		AddString("3", data.EventHash.String()).
		AddString("5", data.InstanceHash.String()).
		AddString("6", data.TaskKey).
		Add("7", data.Inputs).
		AddStringSlice("10", data.Tags).
		AddString("11", data.ProcessHash.String()).
		AddString("12", data.NodeKey).
		AddString("13", data.ExecutorHash.String()).
		AddString("14", data.Price).
		HashSerialize()
}
