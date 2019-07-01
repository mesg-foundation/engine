package executionsdk

import (
	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/x/xstrings"
)

// Filter store fileds for matching executions.
type Filter struct {
	Statuses     []execution.Status
	InstanceHash hash.Hash
	TaskKey      string
	Tags         []string
}

// Match matches execution.
func (f *Filter) Match(e *execution.Execution) bool {
	if f == nil {
		return true
	}

	if !f.InstanceHash.IsZero() && !f.InstanceHash.Equal(e.InstanceHash) {
		return false
	}

	if f.TaskKey != "" && f.TaskKey != "*" && f.TaskKey != e.TaskKey {
		return false
	}

	match := len(f.Statuses) == 0
	for _, status := range f.Statuses {
		if status == e.Status {
			match = true
			break
		}
	}

	if !match {
		return false
	}

	for _, tag := range f.Tags {
		if !xstrings.SliceContains(e.Tags, tag) {
			return false
		}
	}
	return true
}

// HasTaskKey returns true if task key is set to specified value.
func (f *Filter) HasTaskKey() bool {
	return f != nil && f.TaskKey != "" && f.TaskKey != "*"
}

// Listener provides functionalities to listen MESG tasks.
type Listener struct {
	// Channel receives matching executions for tasks.
	C chan *execution.Execution

	ps    *pubsub.PubSub
	topic string
	c     chan interface{}

	filter *Filter
}

// NewListener creates a new Listener with given sdk.
func NewListener(ps *pubsub.PubSub, topic string, f *Filter) *Listener {
	return &Listener{
		C:      make(chan *execution.Execution, 1),
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

// Listen listens executions that match filter.
func (l *Listener) Listen() {
	for v := range l.c {
		if e := v.(*execution.Execution); l.filter.Match(e) {
			l.C <- e
		}
	}
}
