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
	ser := hashserializer.New()
	ser.AddString("2", data.Name)
	ser.Add("4", processNodes(data.Nodes))
	ser.Add("5", processEdges(data.Edges))
	return ser.HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Process_Node) HashSerialize() string {
	if data == nil {
		return ""
	}
	ser := hashserializer.New()
	ser.AddString("1", data.Key)
	ser.Add("2", data.GetResult())
	ser.Add("3", data.GetEvent())
	ser.Add("4", data.GetTask())
	ser.Add("5", data.GetMap())
	ser.Add("6", data.GetFilter())
	return ser.HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Process_Node_Result) HashSerialize() string {
	if data == nil {
		return ""
	}
	ser := hashserializer.New()
	ser.AddString("2", data.InstanceHash.String())
	ser.AddString("3", data.TaskKey)
	return ser.HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Process_Node_Event) HashSerialize() string {
	if data == nil {
		return ""
	}
	ser := hashserializer.New()
	ser.AddString("2", data.InstanceHash.String())
	ser.AddString("3", data.EventKey)
	return ser.HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Process_Node_Task) HashSerialize() string {
	if data == nil {
		return ""
	}
	ser := hashserializer.New()
	ser.AddString("2", data.InstanceHash.String())
	ser.AddString("3", data.TaskKey)
	return ser.HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Process_Node_Map) HashSerialize() string {
	if data == nil {
		return ""
	}
	ser := hashserializer.New()
	ser.Add("1", mapProcessNodeMapOutput(data.Outputs))
	return ser.HashSerialize()
}

type mapProcessNodeMapOutput map[string]*Process_Node_Map_Output

// HashSerialize returns the hashserialized string of this type
func (data mapProcessNodeMapOutput) HashSerialize() string {
	ser := hashserializer.New()
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
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
	ser := hashserializer.New()
	ser.AddInt("1", int(data.GetNull()))
	ser.AddString("2", data.GetStringConst())
	ser.AddFloat("3", data.GetDoubleConst())
	ser.AddBool("4", data.GetBoolConst())
	ser.Add("5", data.GetRef())
	ser.Add("6", data.GetList())
	ser.Add("7", data.GetMap())
	return ser.HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Process_Node_Map_Output_Map) HashSerialize() string {
	if data == nil {
		return ""
	}
	ser := hashserializer.New()
	ser.Add("1", mapProcessNodeMapOutput(data.Outputs))
	return ser.HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Process_Node_Map_Output_List) HashSerialize() string {
	if data == nil {
		return ""
	}
	ser := hashserializer.New()
	ser.Add("1", processNodeMapOutputs(data.Outputs))
	return ser.HashSerialize()
}

type processNodeMapOutputs []*Process_Node_Map_Output

// HashSerialize returns the hashserialized string of this type
func (data processNodeMapOutputs) HashSerialize() string {
	if data == nil {
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
	ser := hashserializer.New()
	ser.AddString("1", data.NodeKey)
	ser.Add("2", data.Path)
	return ser.HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Process_Node_Map_Output_Reference_Path) HashSerialize() string {
	if data == nil {
		return ""
	}
	ser := hashserializer.New()
	ser.AddString("1", data.GetKey())
	ser.AddInt("2", int(data.GetIndex()))
	ser.Add("3", data.Path)
	return ser.HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Process_Node_Filter) HashSerialize() string {
	if data == nil {
		return ""
	}
	ser := hashserializer.New()
	ser.Add("2", processNodeFilterConditions(data.Conditions))
	return ser.HashSerialize()
}

type processNodeFilterConditions []Process_Node_Filter_Condition

// HashSerialize returns the hashserialized string of this type
func (data processNodeFilterConditions) HashSerialize() string {
	if data == nil {
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
	ser := hashserializer.New()
	ser.AddString("1", data.Key)
	ser.AddInt("2", int(data.Predicate))
	ser.AddString("3", data.Value)
	return ser.HashSerialize()
}

// HashSerialize returns the hashserialized string of this type
func (data *Process_Edge) HashSerialize() string {
	if data == nil {
		return ""
	}
	ser := hashserializer.New()
	ser.AddString("1", data.Src)
	ser.AddString("2", data.Dst)
	return ser.HashSerialize()
}

type processNodes []*Process_Node

// HashSerialize returns the hashserialized string of this type
func (data processNodes) HashSerialize() string {
	if data == nil {
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
	if data == nil {
		return ""
	}
	ser := hashserializer.New()
	for i, value := range data {
		ser.Add(strconv.Itoa(i), value)
	}
	return ser.HashSerialize()
}
