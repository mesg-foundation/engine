package api

import (
	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/core/event"
)

// EventFilter store fileds for matching events.
type EventFilter struct {
	Key string
}

// Match matches event.
func (f *EventFilter) Match(e *event.Event) bool {
	return f == nil || f.Key == "" || f.Key == e.Key
}

// EventListener provides functionalities to listen MESG events.
type EventListener struct {
	C chan *event.Event

	ps    *pubsub.PubSub
	topic string
	c     chan interface{}

	filter *EventFilter
}

// NewEventListener creates a new EventListener with given api and filters.
func NewEventListener(ps *pubsub.PubSub, topic string, f *EventFilter) *EventListener {
	return &EventListener{
		C:      make(chan *event.Event, 1),
		ps:     ps,
		topic:  topic,
		c:      ps.Sub(topic),
		filter: f,
	}
}

// Close stops listening for events.
func (l *EventListener) Close() {
	go func() {
		l.ps.Unsub(l.c, l.topic)
		close(l.C)
	}()
}

// Listen listens events that match filter.
func (l *EventListener) Listen() {
	for v := range l.c {
		if e := v.(*event.Event); l.filter.Match(e) {
			l.C <- e
		}
	}
}
