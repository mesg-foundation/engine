package instance

import (
	"github.com/mesg-foundation/engine/hash/serializer"
)

func (data *Instance) Serialize() string {
	if data == nil {
		return ""
	}
	ser := serializer.New()
	ser.AddString("2", data.ServiceHash.String())
	ser.AddString("3", data.EnvHash.String())
	return ser.Serialize()
}
