package ownership

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers the ownership types to codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*isOwnership_Resource)(nil), nil)
	cdc.RegisterConcrete(&Ownership_ServiceHash{}, "mesg.types.Ownership.Ownership_ServiceHash", nil)
}
