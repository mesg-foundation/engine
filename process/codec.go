package process

import (
	"sort"

	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*isProcess_Node_Map_Output_Value)(nil), nil)
	cdc.RegisterConcrete(&Process_Node_Map_Output_Ref{}, "mesg.types.Process_Node_Map_Output_Ref", nil)
	cdc.RegisterConcrete(&Process_Node_Reference{}, "mesg.types.Process_Node_Reference", nil)
	cdc.RegisterConcrete(&Process_Node_Map_Output_Null_{}, "mesg.types.Process_Node_Map_Output_Null_", nil)
	cdc.RegisterConcrete(&Process_Node_Map_Output_StringConst{}, "mesg.types.Process_Node_Map_Output_StringConst", nil)
	cdc.RegisterConcrete(&Process_Node_Map_Output_DoubleConst{}, "mesg.types.Process_Node_Map_Output_DoubleConst", nil)
	cdc.RegisterConcrete(&Process_Node_Map_Output_BoolConst{}, "mesg.types.Process_Node_Map_Output_BoolConst", nil)
	cdc.RegisterConcrete(&Process_Node_Map_Output_List_{}, "mesg.types.Process_Node_Map_Output_List_", nil)
	cdc.RegisterConcrete(&Process_Node_Map_Output_List{}, "mesg.types.Process_Node_Map_Output_List", nil)
	cdc.RegisterConcrete(&Process_Node_Map_Output_Map_{}, "mesg.types.Process_Node_Map_Output_Map_", nil)
	cdc.RegisterConcrete(&Process_Node_Map_Output_Map{}, "mesg.types.Process_Node_Map_Output_Map", nil)
	cdc.RegisterInterface((*isProcess_Node_Type)(nil), nil)
	cdc.RegisterConcrete(&Process_Node_Result_{}, "mesg.types.Process_Node_Result_", nil)
	cdc.RegisterConcrete(&Process_Node_Result{}, "mesg.types.Process_Node_Result", nil)
	cdc.RegisterConcrete(&Process_Node_Event_{}, "mesg.types.Process_Node_Event_", nil)
	cdc.RegisterConcrete(&Process_Node_Event{}, "mesg.types.Process_Node_Event", nil)
	cdc.RegisterConcrete(&Process_Node_Task_{}, "mesg.types.Process_Node_Task_", nil)
	cdc.RegisterConcrete(&Process_Node_Task{}, "mesg.types.Process_Node_Task", nil)
	cdc.RegisterConcrete(&Process_Node_Map_{}, "mesg.types.Process_Node_Map_", nil)
	cdc.RegisterConcrete(&Process_Node_Map{}, "mesg.types.Process_Node_Map", nil)
	cdc.RegisterConcrete(&Process_Node_Filter_{}, "mesg.types.Process_Node_Filter_", nil)
	cdc.RegisterConcrete(&Process_Node_Filter{}, "mesg.types.Process_Node_Filter", nil)
	cdc.RegisterInterface((*isProcess_Node_Reference_Path_Selector)(nil), nil)
	cdc.RegisterConcrete(&Process_Node_Reference_Path_Key{}, "mesg.types.Process_Node_Reference_Path_Key", nil)
	cdc.RegisterConcrete(&Process_Node_Reference_Path_Index{}, "mesg.types.Process_Node_Reference_Path_Index", nil)
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

// MarshalAmino transforms the Process_Node_Map_Output_Map to an array of key/value.
func (m Process_Node_Map_Output_Map) MarshalAmino() ([]KeyOutput, error) {
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

// UnmarshalAmino transforms the key/value array to a Process_Node_Map_Output_Map.
func (m *Process_Node_Map_Output_Map) UnmarshalAmino(keyValues []KeyOutput) error {
	m.Outputs = make(map[string]*Process_Node_Map_Output, len(keyValues))
	for _, p := range keyValues {
		m.Outputs[p.Key] = p.Value
	}
	return nil
}
