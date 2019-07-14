package workflow

import (
	"testing"

	"github.com/mesg-foundation/engine/event"
	"github.com/mesg-foundation/engine/hash"
	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	var tests = []struct {
		trigger *trigger
		event   *event.Event
		match   bool
	}{
		{ // matching event
			&trigger{InstanceHash: hash.Int(1), Key: "xx"},
			&event.Event{InstanceHash: hash.Int(1), Key: "xx"},
			true,
		},
		{ // not matching instance
			&trigger{InstanceHash: hash.Int(1)},
			&event.Event{InstanceHash: hash.Int(2)},
			false,
		},
		{ // not matching event
			&trigger{InstanceHash: hash.Int(1), Key: "xx"},
			&event.Event{InstanceHash: hash.Int(1), Key: "yy"},
			false,
		},
		{ // matching filter
			&trigger{InstanceHash: hash.Int(1), Key: "xx", Filters: []*filter{
				{Key: "foo", Predicate: EQ, Value: "bar"},
			}},
			&event.Event{InstanceHash: hash.Int(1), Key: "xx", Data: map[string]interface{}{"foo": "bar"}},
			true,
		},
		{ // not matching filter
			&trigger{InstanceHash: hash.Int(1), Key: "xx", Filters: []*filter{
				{Key: "foo", Predicate: EQ, Value: "xx"},
			}},
			&event.Event{InstanceHash: hash.Int(1), Key: "xx", Data: map[string]interface{}{"foo": "bar"}},
			false,
		},
		{ // matching multiple filters
			&trigger{InstanceHash: hash.Int(1), Key: "xx", Filters: []*filter{
				{Key: "foo", Predicate: EQ, Value: "bar"},
				{Key: "xxx", Predicate: EQ, Value: "yyy"},
			}},
			&event.Event{InstanceHash: hash.Int(1), Key: "xx", Data: map[string]interface{}{
				"foo": "bar",
				"xxx": "yyy",
				"aaa": "bbb",
			}},
			true,
		},
		{ // non matching multiple filters
			&trigger{InstanceHash: hash.Int(1), Key: "xx", Filters: []*filter{
				{Key: "foo", Predicate: EQ, Value: "bar"},
				{Key: "xxx", Predicate: EQ, Value: "aaa"},
			}},
			&event.Event{InstanceHash: hash.Int(1), Key: "xx", Data: map[string]interface{}{
				"foo": "bar",
				"xxx": "yyy",
				"aaa": "bbb",
			}},
			false,
		},
	}
	for i, test := range tests {
		match := test.trigger.MatchEvent(test.event)
		assert.Equal(t, test.match, match, i)
	}
}
