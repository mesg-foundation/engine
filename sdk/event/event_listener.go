package eventsdk

import (
	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/core/event"
	"github.com/mesg-foundation/core/hash"
)

// Filter store fileds for matching events.
type Filter struct {
	Hash         hash.Hash
	InstanceHash hash.Hash
	Key          string
}

// Match matches event.
func (f *Filter) Match(e *event.Event) bool {
	if f == nil {
		return true
	}

	if !f.Hash.IsZero() && !f.Hash.Equal(e.Hash) {
		return false
	}

	if !f.InstanceHash.IsZero() && !f.InstanceHash.Equal(e.InstanceHash) {
		return false
	}

	if f.Key != "" && f.Key != "*" && f.Key != e.Key {
		return false
	}

	return true
}

// HasKey returns true if key is set to specified value.
func (f *Filter) HasKey() bool {
	return f != nil && f.Key != "" && f.Key != "*"
}

// Listener provides functionalities to listen MESG events.
type Listener struct {
	C chan *event.Event

	ps    *pubsub.PubSub
	topic string
	c     chan interface{}

	filter *Filter
}

// NewListener creates a new Listener with given sdk and filters.
func NewListener(ps *pubsub.PubSub, topic string, f *Filter) *Listener {
	return &Listener{
		C:      make(chan *event.Event, 1),
		ps:     ps,
		topic:  topic,
		c:      ps.Sub(topic),
		filter: f,
	}
}

// Close stops listening for events.
func (l *Listener) Close() {
	go func() {
		l.ps.Unsub(l.c, l.topic)
		close(l.C)
	}()
}

// Listen listens events that match filter.
func (l *Listener) Listen() {
	for v := range l.c {
		if e := v.(*event.Event); l.filter.Match(e) {
			l.C <- e
		}
	}
}
