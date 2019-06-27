package eventsdk

import (
	"testing"

	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/core/event"
	"github.com/mesg-foundation/core/hash"
	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	var tests = []struct {
		f     *Filter
		e     *event.Event
		match bool
	}{
		{
			nil,
			nil,
			true,
		},
		{
			&Filter{},
			&event.Event{},
			true,
		},
		{
			&Filter{Hash: hash.Int(1)},
			&event.Event{Hash: hash.Int(1)},
			true,
		},
		{
			&Filter{Hash: hash.Int(1)},
			&event.Event{Hash: hash.Int(2)},
			false,
		},
		{
			&Filter{InstanceHash: hash.Int(1)},
			&event.Event{InstanceHash: hash.Int(1)},
			true,
		},
		{
			&Filter{InstanceHash: hash.Int(1)},
			&event.Event{InstanceHash: hash.Int(1)},
			true,
		},
		{
			&Filter{Key: "0"},
			&event.Event{Key: "0"},
			true,
		},
		{
			&Filter{Key: "*"},
			&event.Event{Key: "0"},
			true,
		},
		{
			&Filter{Key: "0"},
			&event.Event{Key: "1"},
			false,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.match, tt.f.Match(tt.e))
	}
}

func TestEventListener(t *testing.T) {
	topic := "test-topic"
	testEvent := &event.Event{Key: "0"}
	ps := pubsub.New(0)
	el := NewListener(ps, topic, &Filter{Key: "0"})

	go func() {
		ps.Pub(&event.Event{Key: "1"}, topic)
		ps.Pub(testEvent, topic)
	}()
	go el.Listen()

	recvEvent := <-el.C
	assert.Equal(t, testEvent, recvEvent)

	el.Close()
	_, ok := <-el.C
	assert.False(t, ok)
}
