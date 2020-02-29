package event

import (
	"github.com/mesg-foundation/engine/hash/serializer"
)

func (data *Event) Serialize() string {
	if data == nil {
		return ""
	}
	ser := serializer.New()
	ser.AddString("2", data.InstanceHash.String())
	ser.AddString("3", data.Key)
	ser.Add("4", data.Data)
	return ser.Serialize()
}
