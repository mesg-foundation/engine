package result

import "github.com/mesg-foundation/engine/codec"

func init() {
	codec.RegisterInterface((*isResult_Result)(nil), nil)
	codec.RegisterConcrete(&Result_Outputs{}, "mesg.types.Result.Result_Outputs", nil)
	codec.RegisterConcrete(&Result_Error{}, "mesg.types.Result.Result_Error", nil)
}
