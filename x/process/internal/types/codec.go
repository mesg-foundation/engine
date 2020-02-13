package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	processpb "github.com/mesg-foundation/engine/process"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	processpb.RegisterCodec(cdc)
	cdc.RegisterConcrete(MsgCreateProcess{}, "process/CreateProcess", nil)
	cdc.RegisterConcrete(MsgDeleteProcess{}, "process/DeleteProcess", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
