package ownership

import "github.com/mesg-foundation/engine/codec"

func init() {
	codec.RegisterInterface((*isOwnership_Resource)(nil), nil)
	codec.RegisterConcrete(&Ownership_ServiceHash{}, "mesg.types.Ownership.Ownership_ServiceHash", nil)
}
