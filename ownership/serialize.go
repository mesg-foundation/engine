package ownership

import (
	"github.com/mesg-foundation/engine/hash/serializer"
)

func (data *Ownership) Serialize() string {
	if data == nil {
		return ""
	}
	ser := serializer.New()
	ser.AddString("2", data.Owner)
	ser.AddString("3", data.ResourceHash.String())
	return ser.Serialize()
}
