package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	processpb "github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/protobuf/types"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	processpb.RegisterCodec(cdc)
	cdc.RegisterConcrete(MsgCreate{}, "process/CreateProcess", nil)
	cdc.RegisterConcrete(MsgDelete{}, "process/DeleteProcess", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	types.RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
