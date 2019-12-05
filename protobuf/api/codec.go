package api

import (
	"github.com/mesg-foundation/engine/codec"
)

func init() {
	codec.RegisterInterface((*isCreateResultRequest_Result)(nil), nil)
	codec.RegisterConcrete(&CreateResultRequest_Outputs{}, "mesg.api.CreateResultRequest_Outputs", nil)
	codec.RegisterConcrete(&CreateResultRequest_Error{}, "mesg.api.CreateResultRequest_Error", nil)
}
