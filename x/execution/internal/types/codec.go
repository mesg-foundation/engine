package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/protobuf/types"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreate{}, "execution/create", nil)
	cdc.RegisterConcrete(MsgUpdate{}, "execution/update", nil)
	cdc.RegisterInterface((*isMsgUpdate_Result)(nil), nil)
	cdc.RegisterConcrete(&MsgUpdate_Outputs{}, "mesg.execution.types.MsgUpdate_Outputs", nil)
	cdc.RegisterConcrete(&MsgUpdate_Error{}, "mesg.execution.types.MsgUpdate_Error", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	types.RegisterCodec(ModuleCdc)
	api.RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
