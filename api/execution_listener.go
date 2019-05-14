package api

import (
	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/x/xstrings"
)

// ExecutionFilter store fileds for matching executions.
type ExecutionFilter struct {
	Status    execution.Status
	TaskKey   string
	OutputKey string
	Tags      []string
}

// Match matches execution.
func (f *ExecutionFilter) Match(e *execution.Execution) bool {
	if f == nil {
		return true
	}
	if f.TaskKey != "" && f.TaskKey != "*" && f.TaskKey != e.TaskKey {
		return false
	}
	if f.OutputKey != "" && f.OutputKey != "*" && f.OutputKey != e.OutputKey {
		return false
	}
	if f.Status != 0 && f.Status != e.Status {
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
func (f *ExecutionFilter) HasTaskKey() bool {
	return f != nil && f.TaskKey != "" && f.TaskKey != "*"
}

// ExecutionListener provides functionalities to listen MESG tasks.
type ExecutionListener struct {
	// Channel receives matching executions for tasks.
	C chan *execution.Execution

	ps    *pubsub.PubSub
	topic string
	c     chan interface{}

	filter *ExecutionFilter
}

// NewExecutionListener creates a new ExecutionListener with given api.
func NewExecutionListener(ps *pubsub.PubSub, topic string, f *ExecutionFilter) *ExecutionListener {
	return &ExecutionListener{
		C:      make(chan *execution.Execution, 1),
		ps:     ps,
		topic:  topic,
		c:      ps.Sub(topic),
		filter: f,
	}
}

// Close stops listening for events.
func (l *ExecutionListener) Close() {
	go func() {
		l.ps.Unsub(l.c, l.topic)
		close(l.C)
	}()
}

// Listen listens executions that match filter.
func (l *ExecutionListener) Listen() {
	for v := range l.c {
		if e := v.(*execution.Execution); l.filter.Match(e) {
			l.C <- e
		}
	}
}
