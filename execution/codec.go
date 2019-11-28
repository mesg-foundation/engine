package execution

import "github.com/mesg-foundation/engine/codec"

func init() {
	codec.RegisterInterface((*isExecutionResult_Result)(nil), nil)
	codec.RegisterConcrete(&ExecutionResult_Outputs{}, "mesg.types.Execution.ExecutionResult_Outputs", nil)
	codec.RegisterConcrete(&ExecutionResult_Error{}, "mesg.types.Execution.ExecutionResult_Error", nil)
}
