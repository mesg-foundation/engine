package address

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterCodec registers interface for error.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*sdk.Address)(nil), nil)
	cdc.RegisterConcrete(&EventAddress{}, "mesg.address.EventAddress", nil)
	cdc.RegisterConcrete(&ExecAddress{}, "mesg.address.ExecAddress", nil)
	cdc.RegisterConcrete(&InstAddress{}, "mesg.address.InstAddress", nil)
	cdc.RegisterConcrete(&OwnAddress{}, "mesg.address.OwnAddress", nil)
	cdc.RegisterConcrete(&ProcAddress{}, "mesg.address.ProcAddress", nil)
	cdc.RegisterConcrete(&RunAddress{}, "mesg.address.RunAddress", nil)
	cdc.RegisterConcrete(&ServAddress{}, "mesg.address.ServAddress", nil)
}
