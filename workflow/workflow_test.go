package workflow

import (
	"testing"

	"github.com/mesg-foundation/engine/hash"
	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	var tests = []struct {
		trigger      *Trigger
		instanceHash hash.Hash
		key          string
		data         map[string]interface{}
		match        bool
	}{
		{ // matching event
			trigger:      &Trigger{InstanceHash: hash.Int(1), Key: "xx", Type: EVENT},
			instanceHash: hash.Int(1),
			key:          "xx",
			match:        true,
		},
		{ // not matching instance
			trigger:      &Trigger{InstanceHash: hash.Int(1), Type: EVENT},
			instanceHash: hash.Int(2),
			match:        false,
		},
		{ // not matching event
			trigger:      &Trigger{InstanceHash: hash.Int(1), Key: "xx", Type: EVENT},
			instanceHash: hash.Int(1),
			key:          "yy",
			match:        false,
		},
		{ // matching filter
			trigger: &Trigger{InstanceHash: hash.Int(1), Key: "xx", Type: EVENT, Filters: []*TriggerFilter{
				{Key: "foo", Predicate: EQ, Value: "bar"},
			}},
			instanceHash: hash.Int(1),
			key:          "xx",
			data:         map[string]interface{}{"foo": "bar"},
			match:        true,
		},
		{ // not matching filter
			trigger: &Trigger{InstanceHash: hash.Int(1), Key: "xx", Type: EVENT, Filters: []*TriggerFilter{
				{Key: "foo", Predicate: EQ, Value: "xx"},
			}},
			instanceHash: hash.Int(1),
			key:          "xx",
			data:         map[string]interface{}{"foo": "bar"},
			match:        false,
		},
		{ // matching multiple filters
			trigger: &Trigger{InstanceHash: hash.Int(1), Key: "xx", Type: EVENT, Filters: []*TriggerFilter{
				{Key: "foo", Predicate: EQ, Value: "bar"},
				{Key: "xxx", Predicate: EQ, Value: "yyy"},
			}},
			instanceHash: hash.Int(1),
			key:          "xx",
			data: map[string]interface{}{
				"foo": "bar",
				"xxx": "yyy",
				"aaa": "bbb",
			},
			match: true,
		},
		{ // non matching multiple filters
			trigger: &Trigger{InstanceHash: hash.Int(1), Key: "xx", Type: EVENT, Filters: []*TriggerFilter{
				{Key: "foo", Predicate: EQ, Value: "bar"},
				{Key: "xxx", Predicate: EQ, Value: "aaa"},
			}},
			instanceHash: hash.Int(1),
			key:          "xx",
			data: map[string]interface{}{
				"foo": "bar",
				"xxx": "yyy",
				"aaa": "bbb",
			},
			match: false,
		},
	}
	for i, test := range tests {
		match := test.trigger.Match(EVENT, test.instanceHash, test.key, test.data)
		assert.Equal(t, test.match, match, i)
	}
}

func TestValidateWorkflow(t *testing.T) {

	trigger := Trigger{
		InstanceHash: hash.Int(2),
		Key:          "-",
		Type:         RESULT,
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
		}, err: "workflow should contain edges with one parent maximum"},
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
