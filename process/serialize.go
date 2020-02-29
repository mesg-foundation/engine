package process

import (
	"sort"
	"strconv"

	"github.com/mesg-foundation/engine/hash/serializer"
)

func (data *Process) Serialize() string {
	if data == nil {
		return ""
	}
	ser := serializer.New()
	ser.AddString("2", data.Name)
	ser.Add("4", processNodes(data.Nodes))
	ser.Add("5", processEdges(data.Edges))
	return ser.Serialize()
}

func (data *Process_Node) Serialize() string {
	if data == nil {
		return ""
	}
	ser := serializer.New()
	ser.AddString("1", data.Key)
	ser.Add("2", data.GetResult())
	ser.Add("3", data.GetEvent())
	ser.Add("4", data.GetTask())
	ser.Add("5", data.GetMap())
	ser.Add("6", data.GetFilter())
	return ser.Serialize()
}

func (data *Process_Node_Result) Serialize() string {
	if data == nil {
		return ""
	}
	ser := serializer.New()
	ser.AddString("2", data.InstanceHash.String())
	ser.AddString("3", data.TaskKey)
	return ser.Serialize()
}

func (data *Process_Node_Event) Serialize() string {
	if data == nil {
		return ""
	}
	ser := serializer.New()
	ser.AddString("2", data.InstanceHash.String())
	ser.AddString("3", data.EventKey)
	return ser.Serialize()
}

func (data *Process_Node_Task) Serialize() string {
	if data == nil {
		return ""
	}
	ser := serializer.New()
	ser.AddString("2", data.InstanceHash.String())
	ser.AddString("3", data.TaskKey)
	return ser.Serialize()
}

func (data *Process_Node_Map) Serialize() string {
	if data == nil {
		return ""
	}
	ser := serializer.New()
	ser.Add("1", mapProcessNodeMapOutput(data.Outputs))
	return ser.Serialize()
}

type mapProcessNodeMapOutput map[string]*Process_Node_Map_Output

func (data mapProcessNodeMapOutput) Serialize() string {
	ser := serializer.New()
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		ser.Add(key, data[key])
	}
	return ser.Serialize()
}

func (data *Process_Node_Map_Output) Serialize() string {
	if data == nil {
		return ""
	}
	ser := serializer.New()
	ser.AddString("2", data.GetStringConst())
	ser.AddFloat("3", data.GetDoubleConst())
	ser.AddBool("4", data.GetBoolConst())
	ser.Add("5", data.GetRef())
	ser.Add("6", data.GetList())
	ser.Add("7", data.GetMap())
	return ser.Serialize()
}

func (data *Process_Node_Map_Output_Map) Serialize() string {
	if data == nil {
		return ""
	}
	ser := serializer.New()
	ser.Add("1", mapProcessNodeMapOutput(data.Outputs))
	return ser.Serialize()
}

func (data *Process_Node_Map_Output_List) Serialize() string {
	if data == nil {
		return ""
	}
	ser := serializer.New()
	ser.Add("1", processNodeMapOutputs(data.Outputs))
	return ser.Serialize()
}

type processNodeMapOutputs []*Process_Node_Map_Output

func (data processNodeMapOutputs) Serialize() string {
	if data == nil {
		return ""
	}
	ser := serializer.New()
	for i, value := range data {
		ser.Add(strconv.Itoa(i), value)
	}
	return ser.Serialize()
}

func (data *Process_Node_Map_Output_Reference) Serialize() string {
	if data == nil {
		return ""
	}
	ser := serializer.New()
	ser.AddString("1", data.NodeKey)
	ser.Add("2", data.Path)
	return ser.Serialize()
}

func (data *Process_Node_Map_Output_Reference_Path) Serialize() string {
	if data == nil {
		return ""
	}
	ser := serializer.New()
	ser.AddString("1", data.GetKey())
	ser.AddInt("2", int(data.GetIndex()))
	ser.Add("3", data.Path)
	return ser.Serialize()
}

func (data *Process_Node_Filter) Serialize() string {
	if data == nil {
		return ""
	}
	ser := serializer.New()
	ser.Add("2", processNodeFilterConditions(data.Conditions))
	return ser.Serialize()
}

type processNodeFilterConditions []Process_Node_Filter_Condition

func (data processNodeFilterConditions) Serialize() string {
	if data == nil {
		return ""
	}
	ser := serializer.New()
	for i, value := range data {
		ser.Add(strconv.Itoa(i), value)
	}
	return ser.Serialize()
}

func (data Process_Node_Filter_Condition) Serialize() string {
	ser := serializer.New()
	ser.AddString("1", data.Key)
	ser.AddInt("2", int(data.Predicate))
	return ser.Serialize()
}

func (data *Process_Edge) Serialize() string {
	if data == nil {
		return ""
	}
	ser := serializer.New()
	ser.AddString("1", data.Src)
	ser.AddString("2", data.Dst)
	return ser.Serialize()
}

type processNodes []*Process_Node

func (data processNodes) Serialize() string {
	if data == nil {
		return ""
	}
	ser := serializer.New()
	for i, value := range data {
		ser.Add(strconv.Itoa(i), value)
	}
	return ser.Serialize()
}

type processEdges []*Process_Edge

func (data processEdges) Serialize() string {
	if data == nil {
		return ""
	}
	ser := serializer.New()
	for i, value := range data {
		ser.Add(strconv.Itoa(i), value)
	}
	return ser.Serialize()
}
