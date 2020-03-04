package runner

import "github.com/mesg-foundation/engine/hash/hashserializer"

func (data *Runner) HashSerialize() string {
	if data == nil {
		return ""
	}
	ser := hashserializer.New()
	ser.AddString("2", data.Address)
	ser.AddString("3", data.InstanceHash.String())
	return ser.HashSerialize()
}
