package process

import (
	"encoding/json"
	"fmt"

	"github.com/mesg-foundation/engine/hash"
)

// MarshalJSON for the task
func (t Task) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":         "task",
		"key":          t.Key,
		"instanceHash": t.InstanceHash.String(),
		"taskKey":      t.TaskKey,
	})
}

// MarshalJSON for the result
func (r Result) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":         "result",
		"key":          r.Key,
		"instanceHash": r.InstanceHash.String(),
		"taskKey":      r.TaskKey,
	})
}

// MarshalJSON for the event
func (e Event) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":         "event",
		"key":          e.Key,
		"instanceHash": e.InstanceHash.String(),
		"eventKey":     e.EventKey,
	})
}

// MarshalJSON for the map
func (m Map) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":    "map",
		"key":     m.Key,
		"outputs": m.Outputs,
	})
}

// MarshalJSON for the filter
func (f Filter) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":       "filter",
		"key":        f.Key,
		"conditions": f.Conditions,
	})
}

// UnmarshalJSON unmashals a process
func (w *Process) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*objMap["Hash"], &w.Hash)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*objMap["Key"], &w.Key)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*objMap["Edges"], &w.Edges)
	if err != nil {
		return err
	}

	var rawNodes []*json.RawMessage
	err = json.Unmarshal(*objMap["Nodes"], &rawNodes)
	if err != nil {
		return err
	}

	w.Graph.Nodes = make([]Node, len(rawNodes))
	for i, rawNode := range rawNodes {
		var nodeInfo map[string]interface{}
		err = json.Unmarshal(*rawNode, &nodeInfo)
		if err != nil {
			return err
		}
		data, err := w.preprocessUnmashalNode(nodeInfo)
		if err != nil {
			return err
		}
		marshalData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		w.Graph.Nodes[i], err = w.unmarshalNode(nodeInfo["type"].(string), marshalData)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *Process) unmarshalNode(nodeType string, marshalData []byte) (Node, error) {
	switch nodeType {
	case "task":
		var node Task
		if err := json.Unmarshal(marshalData, &node); err != nil {
			return nil, err
		}
		return &node, nil
	case "event":
		var node Event
		if err := json.Unmarshal(marshalData, &node); err != nil {
			return nil, err
		}
		return &node, nil
	case "result":
		var node Result
		if err := json.Unmarshal(marshalData, &node); err != nil {
			return nil, err
		}
		return &node, nil
	case "map":
		var node Map
		if err := json.Unmarshal(marshalData, &node); err != nil {
			return nil, err
		}
		return &node, nil
	case "filter":
		var node Filter
		if err := json.Unmarshal(marshalData, &node); err != nil {
			return nil, err
		}
		return &node, nil
	default:
		return nil, fmt.Errorf("type %q not supported", nodeType)
	}
}

func (w *Process) preprocessUnmashalNode(nodeInfo map[string]interface{}) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	for key, value := range nodeInfo {
		if key == "type" {
			continue
		}
		if key == "instanceHash" {
			h, err := hash.Decode(value.(string))
			if err != nil {
				return nil, err
			}
			data[key] = h
		} else {
			data[key] = value
		}
	}
	return data, nil
}
