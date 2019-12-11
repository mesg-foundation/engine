package process

import (
	"github.com/mesg-foundation/engine/codec"
	"sort"
)

func init() {
	codec.RegisterInterface((*isProcess_Node_Map_Output_Value)(nil), nil)
	codec.RegisterConcrete(&Process_Node_Map_Output_Ref{}, "mesg.types.Process_Node_Map_Output_Ref", nil)
	codec.RegisterConcrete(&Process_Node_Map_Output_Reference{}, "mesg.types.Process_Node_Map_Output_Reference", nil)
	codec.RegisterConcrete(&Process_Node_Map_Output_Null_{}, "mesg.types.Process_Node_Map_Output_Null_", nil)
	codec.RegisterConcrete(&Process_Node_Map_Output_StringConst{}, "mesg.types.Process_Node_Map_Output_StringConst", nil)
	codec.RegisterConcrete(&Process_Node_Map_Output_DoubleConst{}, "mesg.types.Process_Node_Map_Output_DoubleConst", nil)
	codec.RegisterConcrete(&Process_Node_Map_Output_BoolConst{}, "mesg.types.Process_Node_Map_Output_BoolConst", nil)
	codec.RegisterConcrete(&Process_Node_Map_Output_List_{}, "mesg.types.Process_Node_Map_Output_List_", nil)
	codec.RegisterConcrete(&Process_Node_Map_Output_List{}, "mesg.types.Process_Node_Map_Output_List", nil)
	codec.RegisterConcrete(&Process_Node_Map_Output_Map_{}, "mesg.types.Process_Node_Map_Output_Map_", nil)
	codec.RegisterConcrete(&Process_Node_Map_Output_Map{}, "mesg.types.Process_Node_Map_Output_Map", nil)
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

// KeyOutput is a simple key/value representation of one output of a Process_Node_Map.
type KeyOutput struct {
	Key   string
	Value *Process_Node_Map_Output
}

// MarshalAmino transforms the Process_Node_Map to an array of key/value.
func (m Process_Node_Map) MarshalAmino() ([]KeyOutput, error) {
	p := make([]KeyOutput, len(m.Outputs))
	outputKeys := make([]string, len(m.Outputs))
	i := 0
	for key := range m.Outputs {
		outputKeys[i] = key
		i++
	}
	sort.Stable(sort.StringSlice(outputKeys))
	for i, key := range outputKeys {
		p[i] = KeyOutput{
			Key:   key,
			Value: m.Outputs[key],
		}
	}
	return p, nil
}

// UnmarshalAmino transforms the key/value array to a Process_Node_Map.
func (m *Process_Node_Map) UnmarshalAmino(keyValues []KeyOutput) error {
	m.Outputs = make(map[string]*Process_Node_Map_Output, len(keyValues))
	for _, p := range keyValues {
		m.Outputs[p.Key] = p.Value
	}
	return nil
}
