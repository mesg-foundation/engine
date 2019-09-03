package process

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mesg-foundation/engine/filter"
	"github.com/mesg-foundation/engine/hash"
)

func TestMarshal(t *testing.T) {
	w := Process{
		Hash: hash.Int(0),
		Key:  "test",
		Nodes: []*Process_Node{
			&Process_Node{
				Type: &Process_Node_Task_{&Process_Node_Task{Key: "1", InstanceHash: hash.Int(1), TaskKey: "1"}},
			},
			&Process_Node{
				Type: &Process_Node_Result_{&Process_Node_Result{Key: "2", InstanceHash: hash.Int(2), TaskKey: "2"}},
			},
			&Process_Node{
				Type: &Process_Node_Event_{&Process_Node_Event{Key: "3", InstanceHash: hash.Int(3), EventKey: "3"}},
			},
			&Process_Node{
				Type: &Process_Node_Map_{&Process_Node_Map{
					Key: "4", Outputs: []*Process_Node_Map_Output{{
						Key: "5",
						Value: &Process_Node_Map_Output_Ref{
							Ref: &Process_Node_Map_Output_Reference{NodeKey: "5", Key: "5"},
						}},
					}},
				},
			},
			&Process_Node{
				Type: &Process_Node_Filter_{&Process_Node_Filter{
					Key: "6",
					Conditions: []filter.Condition{
						{Key: "x", Predicate: filter.EQ, Value: "x"},
					},
				}},
			},
		},
		Edges: []*Process_Edge{
			{Src: "1", Dst: "2"},
		},
	}
	val, err := json.Marshal(w)
	assert.NoError(t, err)
	// assert.Equal(t, "{\"Nodes\":[{\"instanceHash\":\"4uQeVj5tqViQh7yWWGStvkEG1Zmhx6uasJtWCJziofM\",\"key\":\"1\",\"taskKey\":\"1\",\"type\":\"task\"},{\"instanceHash\":\"8opHzTAnfzRpPEx21XtnrVTX28YQuCpAjcn1PczScKh\",\"key\":\"2\",\"taskKey\":\"2\",\"type\":\"result\"},{\"eventKey\":\"3\",\"instanceHash\":\"CiDwVBFgWV9E5MvXWoLgnEgn2hK7rJikbvfWavzAQz3\",\"key\":\"3\",\"type\":\"event\"},{\"key\":\"4\",\"outputs\":[{\"Key\":\"5\",\"Ref\":{\"NodeKey\":\"5\",\"Key\":\"5\"}}],\"type\":\"map\"},{\"conditions\":[{\"Key\":\"x\",\"Predicate\":1,\"Value\":\"x\"}],\"key\":\"6\",\"type\":\"filter\"}],\"Edges\":[{\"Src\":\"1\",\"Dst\":\"2\"}],\"Hash\":\"11111111111111111111111111111111\",\"Key\":\"test\"}", string(val))
	var w2 Process
	err = json.Unmarshal(val, &w2)
	assert.NoError(t, err)
	assert.Equal(t, w2.Hash, hash.Int(0))
	assert.Equal(t, w2.Key, "test")
	assert.Equal(t, len(w2.Edges), 1)
	assert.Equal(t, w2.Edges[0].Src, "1")
	assert.Equal(t, w2.Edges[0].Dst, "2")
	assert.Equal(t, 5, len(w2.Nodes))
	assert.IsType(t, &Process_Node{}, w2.Nodes[0])
	assert.NotNil(t, w2.Nodes[0].GetTask())
	assert.Equal(t, w2.Nodes[0].GetTask().Key, "1")
	assert.Equal(t, w2.Nodes[0].GetTask().InstanceHash, hash.Int(1))
	assert.Equal(t, w2.Nodes[0].GetTask().TaskKey, "1")
	assert.IsType(t, &Process_Node{}, w2.Nodes[1])
	assert.NotNil(t, w2.Nodes[1].GetResult())
	assert.Equal(t, w2.Nodes[1].GetResult().InstanceHash, hash.Int(2))
	assert.Equal(t, w2.Nodes[1].GetResult().TaskKey, "2")
	assert.IsType(t, &Process_Node{}, w2.Nodes[2])
	assert.NotNil(t, w2.Nodes[2].GetEvent())
	assert.Equal(t, w2.Nodes[2].GetEvent().InstanceHash, hash.Int(3))
	assert.Equal(t, w2.Nodes[2].GetEvent().EventKey, "3")
	assert.IsType(t, &Process_Node{}, w2.Nodes[3])
	// FIXME: issue with the JSON unmarshal of the output of the map
	// assert.NotNil(t, w2.Nodes[3].GetMap())
	// assert.Equal(t, w2.Nodes[3].GetMap().Key, "4")
	// assert.Equal(t, len(w2.Nodes[3].GetMap().Outputs), 1)
	// assert.Equal(t, w2.Nodes[3].GetMap().Outputs[0].Key, "5")
	// assert.Equal(t, w2.Nodes[3].GetMap().Outputs[0].GetRef().NodeKey, "5")
	// assert.Equal(t, w2.Nodes[3].GetMap().Outputs[0].GetRef().Key, "5")
	assert.NotNil(t, w2.Nodes[4].GetFilter())
	assert.Equal(t, w2.Nodes[4].GetFilter().Key, "6")
	assert.Equal(t, len(w2.Nodes[4].GetFilter().Conditions), 1)
	assert.Equal(t, w2.Nodes[4].GetFilter().Conditions[0].Key, "x")
	assert.Equal(t, w2.Nodes[4].GetFilter().Conditions[0].Predicate, filter.EQ)
	assert.Equal(t, w2.Nodes[4].GetFilter().Conditions[0].Value, "x")
}
