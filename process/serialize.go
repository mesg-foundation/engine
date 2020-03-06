package process

import (
	"sort"
	"strconv"

	"github.com/mesg-foundation/engine/hash/hashserializer"
)

// HashSerialize returns the hashserialized string of this type
func (data *Process) HashSerialize() string {
	if data == nil {
		return ""
	}
	return hashserializer.New().
		AddString("2", data.Name).
		Add("4", processNodes(data.Nodes)).
		Add("5", processEdges(data.Edges)).
		HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Process_Node) HashSerialize() string {
	if data == nil {
		return ""
	}
	return hashserializer.New().
		AddString("1", data.Key).
		Add("2", data.GetResult()).
		Add("3", data.GetEvent()).
		Add("4", data.GetTask()).
		Add("5", data.GetMap()).
		Add("6", data.GetFilter()).
		HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Process_Node_Result) HashSerialize() string {
	if data == nil {
		return ""
	}
	return hashserializer.New().
		AddString("2", data.InstanceHash.String()).
		AddString("3", data.TaskKey).
		HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Process_Node_Event) HashSerialize() string {
	if data == nil {
		return ""
	}
	return hashserializer.New().
		AddString("2", data.InstanceHash.String()).
		AddString("3", data.EventKey).
		HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Process_Node_Task) HashSerialize() string {
	if data == nil {
		return ""
	}
	return hashserializer.New().
		AddString("2", data.InstanceHash.String()).
		AddString("3", data.TaskKey).
		HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Process_Node_Map) HashSerialize() string {
	if data == nil {
		return ""
	}
	return hashserializer.New().
		Add("1", mapProcessNodeMapOutput(data.Outputs)).
		HashSerialize()
}

type mapProcessNodeMapOutput map[string]*Process_Node_Map_Output

// HashSerialize returns the hashserialized string of this type
func (data mapProcessNodeMapOutput) HashSerialize() string {
	if data == nil || len(data) == 0 {
		return ""
	}
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	ser := hashserializer.New()
	for _, key := range keys {
		ser.Add(key, data[key])
	}
	return ser.HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Process_Node_Map_Output) HashSerialize() string {
	if data == nil {
		return ""
	}
	return hashserializer.New().
		AddInt("1", int(data.GetNull())).
		AddString("2", data.GetStringConst()).
		AddFloat("3", data.GetDoubleConst()).
		AddBool("4", data.GetBoolConst()).
		Add("5", data.GetRef()).
		Add("6", data.GetList()).
		Add("7", data.GetMap()).
		HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Process_Node_Map_Output_Map) HashSerialize() string {
	if data == nil {
		return ""
	}
	return hashserializer.New().
		Add("1", mapProcessNodeMapOutput(data.Outputs)).
		HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Process_Node_Map_Output_List) HashSerialize() string {
	if data == nil {
		return ""
	}
	return hashserializer.New().
		Add("1", processNodeMapOutputs(data.Outputs)).
		HashSerialize()
}

type processNodeMapOutputs []*Process_Node_Map_Output

// HashSerialize returns the hashserialized string of this type
func (data processNodeMapOutputs) HashSerialize() string {
	if data == nil || len(data) == 0 {
		return ""
	}
	ser := hashserializer.New()
	for i, value := range data {
		ser.Add(strconv.Itoa(i), value)
	}
	return ser.HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Process_Node_Map_Output_Reference) HashSerialize() string {
	if data == nil {
		return ""
	}
	return hashserializer.New().
		AddString("1", data.NodeKey).
		Add("2", data.Path).
		HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Process_Node_Map_Output_Reference_Path) HashSerialize() string {
	if data == nil {
		return ""
	}
	return hashserializer.New().
		AddString("1", data.GetKey()).
		AddInt("2", int(data.GetIndex())).
		Add("3", data.Path).
		HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Process_Node_Filter) HashSerialize() string {
	if data == nil {
		return ""
	}
	return hashserializer.New().
		Add("2", processNodeFilterConditions(data.Conditions)).
		HashSerialize()
}

type processNodeFilterConditions []Process_Node_Filter_Condition

// HashSerialize returns the hashserialized string of this type
func (data processNodeFilterConditions) HashSerialize() string {
	if data == nil || len(data) == 0 {
		return ""
	}
	ser := hashserializer.New()
	for i, value := range data {
		ser.Add(strconv.Itoa(i), value)
	}
	return ser.HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data Process_Node_Filter_Condition) HashSerialize() string {
	return hashserializer.New().
		AddString("1", data.Key).
		AddInt("2", int(data.Predicate)).
		AddString("3", data.Value).
		HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Process_Edge) HashSerialize() string {
	if data == nil {
		return ""
	}
	return hashserializer.New().
		AddString("1", data.Src).
		AddString("2", data.Dst).
		HashSerialize()
}

type processNodes []*Process_Node

// HashSerialize returns the hashserialized string of this type
func (data processNodes) HashSerialize() string {
	if data == nil || len(data) == 0 {
		return ""
	}
	ser := hashserializer.New()
	for i, value := range data {
		ser.Add(strconv.Itoa(i), value)
	}
	return ser.HashSerialize()
}

type processEdges []*Process_Edge

// HashSerialize returns the hashserialized string of this type
func (data processEdges) HashSerialize() string {
	if data == nil || len(data) == 0 {
		return ""
	}
	ser := hashserializer.New()
	for i, value := range data {
		ser.Add(strconv.Itoa(i), value)
	}
	return ser.HashSerialize()
}
