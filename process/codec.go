package process

import (
	"github.com/mesg-foundation/engine/codec"
)

func init() {
	codec.RegisterInterface((*isProcess_Node_Map_Output_Value)(nil), nil)
	codec.RegisterConcrete(&Process_Node_Map_Output_Ref{}, "mesg.types.Process_Node_Map_Output_Ref", nil)
	codec.RegisterConcrete(&Process_Node_Map_Output_Constant{}, "mesg.types.Process_Node_Map_Output_Constant", nil)
	codec.RegisterInterface((*isProcess_Node_Type)(nil), nil)
	codec.RegisterConcrete(&Process_Node_Result_{}, "mesg.types.Process_Node_Result_", nil)
	codec.RegisterConcrete(&Process_Node_Result{}, "mesg.types.Process_Node_Result", nil)
	codec.RegisterConcrete(&Process_Node_Event_{}, "mesg.types.Process_Node_Event_", nil)
	codec.RegisterConcrete(&Process_Node_Event{}, "mesg.types.Process_Node_Event", nil)
	codec.RegisterConcrete(&Process_Node_Task_{}, "mesg.types.Process_Node_Task_", nil)
	codec.RegisterConcrete(&Process_Node_Task{}, "mesg.types.Process_Node_Task", nil)
	codec.RegisterConcrete(&Process_Node_Map_{}, "mesg.types.Process_Node_Map_", nil)
	codec.RegisterConcrete(&Process_Node_Map{}, "mesg.types.Process_Node_Map", nil)
	codec.RegisterConcrete(&Process_Node_Filter_{}, "mesg.types.Process_Node_Filter_", nil)
	codec.RegisterConcrete(&Process_Node_Filter{}, "mesg.types.Process_Node_Filter", nil)
}
