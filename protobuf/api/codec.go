package api

import (
	"github.com/mesg-foundation/engine/codec"
)

func init() {
	codec.RegisterInterface((*isUpdateExecutionRequest_Result)(nil), nil)
	codec.RegisterConcrete(&UpdateExecutionRequest_Outputs{}, "mesg.api.UpdateExecutionRequest_Outputs", nil)
	codec.RegisterConcrete(&UpdateExecutionRequest_Error{}, "mesg.api.UpdateExecutionRequest_Error", nil)
}
