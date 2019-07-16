package service

import (
	"testing"

	"github.com/mesg-foundation/engine/hash"
	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	var tests = []struct {
		trigger      *WorkflowTrigger
		instanceHash hash.Hash
		key          string
		data         map[string]interface{}
		match        bool
	}{
		{ // matching event
			trigger:      &WorkflowTrigger{InstanceHash: hash.Int(1), Key: "xx"},
			instanceHash: hash.Int(1),
			key:          "xx",
			match:        true,
		},
		{ // not matching instance
			trigger:      &WorkflowTrigger{InstanceHash: hash.Int(1)},
			instanceHash: hash.Int(2),
			match:        false,
		},
		{ // not matching event
			trigger:      &WorkflowTrigger{InstanceHash: hash.Int(1), Key: "xx"},
			instanceHash: hash.Int(1),
			key:          "yy",
			match:        false,
		},
		{ // matching filter
			trigger: &WorkflowTrigger{InstanceHash: hash.Int(1), Key: "xx", Filters: []*WorkflowTriggerFilter{
				{Key: "foo", Predicate: EQ, Value: "bar"},
			}},
			instanceHash: hash.Int(1),
			key:          "xx",
			data:         map[string]interface{}{"foo": "bar"},
			match:        true,
		},
		{ // not matching filter
			trigger: &WorkflowTrigger{InstanceHash: hash.Int(1), Key: "xx", Filters: []*WorkflowTriggerFilter{
				{Key: "foo", Predicate: EQ, Value: "xx"},
			}},
			instanceHash: hash.Int(1),
			key:          "xx",
			data:         map[string]interface{}{"foo": "bar"},
			match:        false,
		},
		{ // matching multiple filters
			trigger: &WorkflowTrigger{InstanceHash: hash.Int(1), Key: "xx", Filters: []*WorkflowTriggerFilter{
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
			trigger: &WorkflowTrigger{InstanceHash: hash.Int(1), Key: "xx", Filters: []*WorkflowTriggerFilter{
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
