package runner

import (
	"github.com/mesg-foundation/engine/hash/serializer"
)

func (data *Runner) Serialize() string {
	if data == nil {
		return ""
	}
	ser := serializer.New()
	ser.AddString("2", data.Address)
	ser.AddString("3", data.InstanceHash.String())
	return ser.Serialize()
}
