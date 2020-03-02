package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/protobuf/types"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateExecution{}, "execution/CreateExecution", nil)
	cdc.RegisterConcrete(MsgUpdateExecution{}, "execution/UpdateExecution", nil)
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
