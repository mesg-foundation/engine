package workflow

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mesg-foundation/engine/filter"
	"github.com/mesg-foundation/engine/hash"
)

func TestMarshal(t *testing.T) {
	w := Workflow{
		Hash: hash.Int(0),
		Key:  "test",
		Graph: Graph{
			Nodes: []Node{
				Task{Key: "1", InstanceHash: hash.Int(1), TaskKey: "1"},
				Result{Key: "2", InstanceHash: hash.Int(2), TaskKey: "2"},
				Event{Key: "3", InstanceHash: hash.Int(3), EventKey: "3"},
				Map{Key: "4", Outputs: []Output{
					{Key: "5", Ref: &OutputReference{NodeKey: "5", Key: "5"}},
				}},
				Filter{Key: "6", Filter: filter.Filter{Conditions: []filter.Condition{
					{Key: "7", Predicate: filter.EQ, Value: "7"},
				}}},
			},
			Edges: []Edge{
				{Src: "1", Dst: "2"},
			},
		},
	}
	val, err := json.Marshal(w)
	assert.NoError(t, err)
	assert.Equal(t, "{\"Nodes\":[{\"instanceHash\":\"4uQeVj5tqViQh7yWWGStvkEG1Zmhx6uasJtWCJziofM\",\"key\":\"1\",\"taskKey\":\"1\",\"type\":\"task\"},{\"instanceHash\":\"8opHzTAnfzRpPEx21XtnrVTX28YQuCpAjcn1PczScKh\",\"taskKey\":\"2\",\"type\":\"result\"},{\"eventKey\":\"3\",\"instanceHash\":\"CiDwVBFgWV9E5MvXWoLgnEgn2hK7rJikbvfWavzAQz3\",\"type\":\"event\"},{\"key\":\"4\",\"outputs\":[{\"Key\":\"5\",\"Ref\":{\"NodeKey\":\"5\",\"Key\":\"5\"}}],\"type\":\"mapping\"},{\"conditions\":[{\"Key\":\"7\",\"Predicate\":1,\"Value\":\"7\"}],\"key\":\"6\",\"type\":\"filter\"}],\"Edges\":[{\"Src\":\"1\",\"Dst\":\"2\"}],\"Hash\":\"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=\",\"Key\":\"test\"}", string(val))
	var w2 Workflow
	err = json.Unmarshal(val, &w2)
	assert.NoError(t, err)
	assert.Equal(t, w2.Hash, hash.Int(0))
	assert.Equal(t, w2.Key, "test")
	assert.Equal(t, len(w2.Edges), 1)
	assert.Equal(t, w2.Edges[0].Src, "1")
	assert.Equal(t, w2.Edges[0].Dst, "2")
	assert.Equal(t, 5, len(w2.Nodes))
	assert.IsType(t, &Task{}, w2.Nodes[0])
	assert.Equal(t, w2.Nodes[0].(*Task).Key, "1")
	assert.Equal(t, w2.Nodes[0].(*Task).InstanceHash, hash.Int(1))
	assert.Equal(t, w2.Nodes[0].(*Task).TaskKey, "1")
	assert.IsType(t, &Result{}, w2.Nodes[1])
	assert.Equal(t, w2.Nodes[1].(*Result).InstanceHash, hash.Int(2))
	assert.Equal(t, w2.Nodes[1].(*Result).TaskKey, "2")
	assert.IsType(t, &Event{}, w2.Nodes[2])
	assert.Equal(t, w2.Nodes[2].(*Event).InstanceHash, hash.Int(3))
	assert.Equal(t, w2.Nodes[2].(*Event).EventKey, "3")
	assert.IsType(t, &Map{}, w2.Nodes[3])
	assert.Equal(t, w2.Nodes[3].(*Map).Key, "4")
	assert.Equal(t, len(w2.Nodes[3].(*Map).Outputs), 1)
	assert.Equal(t, w2.Nodes[3].(*Map).Outputs[0].Key, "5")
	assert.Equal(t, w2.Nodes[3].(*Map).Outputs[0].Ref.NodeKey, "5")
	assert.Equal(t, w2.Nodes[3].(*Map).Outputs[0].Ref.Key, "5")
	assert.Equal(t, w2.Nodes[4].(*Filter).Key, "6")
	assert.Equal(t, w2.Nodes[4].(*Filter).Conditions[0].Key, "7")
	assert.Equal(t, w2.Nodes[4].(*Filter).Conditions[0].Predicate, filter.EQ)
	assert.Equal(t, w2.Nodes[4].(*Filter).Conditions[0].Value, "7")
}
