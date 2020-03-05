package event

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cskr/pubsub"
)

// Filter store fileds for matching events.
type Filter struct {
	Hash         sdk.AccAddress
	InstanceHash sdk.AccAddress
	Key          string
}

// Match matches event.
func (f *Filter) Match(e *Event) bool {
	if f == nil {
		return true
	}

	if !f.Hash.Empty() && !f.Hash.Equals(e.Hash) {
		return false
	}

	if !f.InstanceHash.Empty() && !f.InstanceHash.Equals(e.InstanceHash) {
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
	C chan *Event

	ps    *pubsub.PubSub
	topic string
	c     chan interface{}

	filter *Filter
}

// NewListener creates a new Listener with given sdk and filters.
func NewListener(ps *pubsub.PubSub, topic string, f *Filter) *Listener {
	return &Listener{
		C:      make(chan *Event, 1),
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
		if e := v.(*Event); l.filter.Match(e) {
			l.C <- e
		}
	}
}
