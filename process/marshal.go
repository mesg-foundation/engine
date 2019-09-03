package process

import (
	"encoding/json"
	"fmt"

	"github.com/mesg-foundation/engine/hash"
)


// UnmarshalJSON unmashals a process
func (w *Process) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*objMap["hash"], &w.Hash)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*objMap["key"], &w.Key)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*objMap["edges"], &w.Edges)
	if err != nil {
		return err
	}

	var rawNodes []*json.RawMessage
	err = json.Unmarshal(*objMap["nodes"], &rawNodes)
	if err != nil {
		return err
	}

	w.Nodes = make([]*Process_Node, len(rawNodes))
	for i, rawNode := range rawNodes {
		var nodeInfo map[string]interface{}
		err = json.Unmarshal(*rawNode, &nodeInfo)
		if err != nil {
			return err
		}

		for nodeType, data := range nodeInfo["Type"].(map[string]interface{}) {
			marshalData, err := json.Marshal(data)
			if err != nil {
				return err
			}
			w.Nodes[i], err = w.unmarshalNode(nodeType, marshalData)
		}
	}
	return nil
}

func (w *Process) unmarshalNode(nodeType string, marshalData []byte) (*Process_Node, error) {
	switch nodeType {
	case "Task":
		var node Process_Node_Task
		if err := json.Unmarshal(marshalData, &node); err != nil {
			return nil, err
		}
		return &Process_Node{Type: &Process_Node_Task_{Task: &node}}, nil
	case "Event":
		var node Process_Node_Event
		if err := json.Unmarshal(marshalData, &node); err != nil {
			return nil, err
		}
		return &Process_Node{Type: &Process_Node_Event_{Event: &node}}, nil
	case "Result":
		var node Process_Node_Result
		if err := json.Unmarshal(marshalData, &node); err != nil {
			return nil, err
		}
		return &Process_Node{Type: &Process_Node_Result_{Result: &node}}, nil
	case "Map":
		var node Process_Node_Map
		if err := json.Unmarshal(marshalData, &node); err != nil {
			return nil, err
		}
		return &Process_Node{Type: &Process_Node_Map_{Map: &node}}, nil
	case "Filter":
		var node Process_Node_Filter
		if err := json.Unmarshal(marshalData, &node); err != nil {
			return nil, err
		}
		return &Process_Node{Type: &Process_Node_Filter_{Filter: &node}}, nil
	default:
		return nil, fmt.Errorf("type %q not supported", nodeType)
	}
}
