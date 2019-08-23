package workflow

import (
	"testing"

	"github.com/mesg-foundation/engine/hash"
	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	var tests = []struct {
		filters TriggerFilters
		data    map[string]interface{}
		match   bool
	}{
		{ // not matching filter
			filters: []*TriggerFilter{
				{Key: "foo", Predicate: EQ, Value: "xx"},
			},
			data:  map[string]interface{}{"foo": "bar"},
			match: false,
		},
		{ // matching multiple filters
			filters: []*TriggerFilter{
				{Key: "foo", Predicate: EQ, Value: "bar"},
				{Key: "xxx", Predicate: EQ, Value: "yyy"},
			},
			data: map[string]interface{}{
				"foo": "bar",
				"xxx": "yyy",
				"aaa": "bbb",
			},
			match: true,
		},
		{ // non matching multiple filters
			filters: []*TriggerFilter{
				{Key: "foo", Predicate: EQ, Value: "bar"},
				{Key: "xxx", Predicate: EQ, Value: "aaa"},
			},
			data: map[string]interface{}{
				"foo": "bar",
				"xxx": "yyy",
				"aaa": "bbb",
			},
			match: false,
		},
	}
	for i, test := range tests {
		match := test.filters.Match(test.data)
		assert.Equal(t, test.match, match, i)
	}
}

func TestValidateWorkflow(t *testing.T) {

	trigger := Trigger{
		InstanceHash: hash.Int(2),
		TaskKey:      "-",
		NodeKey:      "nodeKey1",
	}

	nodes := []Node{
		{
			Key:          "nodeKey1",
			InstanceHash: hash.Int(2),
			TaskKey:      "-",
		},
		{
			Key:          "nodeKey2",
			InstanceHash: hash.Int(3),
			TaskKey:      "-",
		},
	}

	var tests = []struct {
		w     *Workflow
		valid bool
		err   string
	}{
		{w: &Workflow{
			Hash: hash.Int(1),
			Key:  "invalid-struct",
		}, err: "Error:Field validation"},
		{w: &Workflow{
			Trigger: Trigger{InstanceHash: hash.Int(1), NodeKey: "-"},
			Hash:    hash.Int(1),
			Key:     "missing-key",
		}, err: "Key: 'Workflow.Trigger.TaskKey' Error:Field validation for 'TaskKey' failed on the 'required_without' tag"},
		{w: &Workflow{
			Hash:    hash.Int(1),
			Key:     "trigger-missing-node",
			Trigger: trigger,
		}, err: "node \"nodeKey1\" not found"},
		{w: &Workflow{
			Hash:    hash.Int(1),
			Key:     "edge-src-missing-node",
			Trigger: trigger,
			Nodes:   nodes,
			Edges: []Edge{
				{Src: "-", Dst: "nodeKey2"},
			},
		}, err: "node \"-\" not found"},
		{w: &Workflow{
			Hash:    hash.Int(1),
			Key:     "edge-dst-missing-node",
			Trigger: trigger,
			Nodes:   nodes,
			Edges: []Edge{
				{Src: "nodeKey1", Dst: "-"},
			},
		}, err: "node \"-\" not found"},
		{w: &Workflow{
			Hash:    hash.Int(1),
			Key:     "cyclic-graph",
			Trigger: trigger,
			Nodes:   nodes,
			Edges: []Edge{
				{Src: "nodeKey1", Dst: "nodeKey2"},
				{Src: "nodeKey2", Dst: "nodeKey1"},
			},
		}, err: "workflow should not contain any cycles"},
		{w: &Workflow{
			Hash:    hash.Int(1),
			Key:     "non-connected-graph",
			Trigger: trigger,
			Nodes: append(nodes, Node{
				Key:          "nodeKey3",
				InstanceHash: hash.Int(2),
				TaskKey:      "-",
			}, Node{
				Key:          "nodeKey4",
				InstanceHash: hash.Int(2),
				TaskKey:      "-",
			}),
			Edges: []Edge{
				{Src: "nodeKey1", Dst: "nodeKey2"},
				{Src: "nodeKey3", Dst: "nodeKey4"},
			},
		}, err: "workflow should be a connected graph"},
		{w: &Workflow{
			Hash:    hash.Int(1),
			Key:     "multiple-parent-graph",
			Trigger: trigger,
			Nodes: append(nodes, Node{
				Key:          "nodeKey3",
				InstanceHash: hash.Int(2),
				TaskKey:      "-",
			}, Node{
				Key:          "nodeKey4",
				InstanceHash: hash.Int(2),
				TaskKey:      "-",
			}),
			Edges: []Edge{
				{Src: "nodeKey1", Dst: "nodeKey2"},
				{Src: "nodeKey1", Dst: "nodeKey3"},
				{Src: "nodeKey2", Dst: "nodeKey4"},
				{Src: "nodeKey3", Dst: "nodeKey4"},
			},
		}, err: "workflow should contain nodes with one parent maximum"},
		{w: &Workflow{
			Hash:    hash.Int(1),
			Key:     "multiple-parent-graph",
			Trigger: trigger,
			Nodes: append(nodes, Node{
				Key:          "nodeKey3",
				InstanceHash: hash.Int(2),
				TaskKey:      "-",
			}, Node{
				Key:          "nodeKey4",
				InstanceHash: hash.Int(2),
				TaskKey:      "-",
			}, Node{
				Key:          "nodeKey5",
				InstanceHash: hash.Int(2),
				TaskKey:      "-",
			}, Node{
				Key:          "nodeKey6",
				InstanceHash: hash.Int(2),
				TaskKey:      "-",
			}, Node{
				Key:          "nodeKey7",
				InstanceHash: hash.Int(2),
				TaskKey:      "-",
			}),
			Edges: []Edge{
				{Src: "nodeKey1", Dst: "nodeKey2"},
				{Src: "nodeKey2", Dst: "nodeKey3"},
				{Src: "nodeKey2", Dst: "nodeKey4"},
				{Src: "nodeKey3", Dst: "nodeKey5"},
				{Src: "nodeKey4", Dst: "nodeKey6"},
				{Src: "nodeKey4", Dst: "nodeKey7"},
			},
		}, valid: true},
		{w: &Workflow{
			Hash:    hash.Int(1),
			Key:     "inputs-with-invalid-node",
			Trigger: trigger,
			Nodes:   nodes,
			Edges: []Edge{
				{Src: "nodeKey1", Dst: "nodeKey2", Inputs: []*Input{
					{Key: "-", Ref: &InputReference{Key: "-", NodeKey: "invalid"}},
				}},
			},
		}, err: "node \"invalid\" not found"},
		{w: &Workflow{
			Hash:    hash.Int(1),
			Key:     "inputs-with-valid-ref",
			Trigger: trigger,
			Nodes:   nodes,
			Edges: []Edge{
				{Src: "nodeKey1", Dst: "nodeKey2", Inputs: []*Input{
					{Key: "-", Ref: &InputReference{Key: "-", NodeKey: "nodeKey1"}},
				}},
			},
		}, valid: true},
	}
	for _, test := range tests {
		err := test.w.Validate()
		if test.valid {
			assert.Nil(t, err, test.w.Key)
		} else {
			assert.Contains(t, test.w.Validate().Error(), test.err, test.w.Key)
		}
	}
}
