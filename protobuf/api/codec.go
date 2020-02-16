package api

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers interface for error.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*isUpdateExecutionRequest_Result)(nil), nil)
	cdc.RegisterConcrete(&UpdateExecutionRequest_Outputs{}, "mesg.api.UpdateExecutionRequest_Outputs", nil)
	cdc.RegisterConcrete(&UpdateExecutionRequest_Error{}, "mesg.api.UpdateExecutionRequest_Error", nil)
}
