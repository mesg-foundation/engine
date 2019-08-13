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

	edges := []Edge{
		{
			Src: "nodeKey1",
			Dst: "nodeKey2",
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
			Key:     "trigger-valid-node",
			Trigger: trigger,
			Nodes:   nodes,
		}, valid: true},
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
			Key:     "valid-edges",
			Trigger: trigger,
			Nodes:   nodes,
			Edges:   edges,
		}, valid: true},
	}
	for _, test := range tests {
		err := test.w.Validate()
		if test.valid {
			assert.Nil(t, err)
		} else {
			assert.Contains(t, test.w.Validate().Error(), test.err)
		}
	}
}
