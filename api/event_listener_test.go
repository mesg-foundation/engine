package api

import (
	"testing"

	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/core/event"
	"github.com/stretchr/testify/assert"
)

func TestEventFilter(t *testing.T) {
	var tests = []struct {
		f     *EventFilter
		e     *event.Event
		match bool
	}{
		{
			nil,
			nil,
			true,
		},
		{
			&EventFilter{},
			&event.Event{},
			true,
		},
		{
			&EventFilter{Key: "0"},
			&event.Event{Key: "0"},
			true,
		},
		{
			&EventFilter{Key: "0"},
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
	el := NewEventListener(ps, topic, &EventFilter{Key: "0"})

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
