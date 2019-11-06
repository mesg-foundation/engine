package runnersdk

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc is the codec for the module
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	ModuleCdc.RegisterConcrete(msgCreateRunner{}, "runner/create", nil)
	ModuleCdc.RegisterConcrete(msgDeleteRunner{}, "runner/delete", nil)
}
