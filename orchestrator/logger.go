package orchestrator

import (
	"github.com/cskr/pubsub"
)

// Logger provides functionalities to listen to the orchestrator's logs
type Logger struct {
	Logs chan *OrchestratorLog

	ps     *pubsub.PubSub
	topics []string
	sub    chan interface{}
}

// NewLogger creates a new logger for the given topics
func (s *Orchestrator) NewLogger(topics ...string) *Logger {
	return &Logger{
		Logs:   make(chan *OrchestratorLog),
		ps:     s.logs,
		topics: topics,
		sub:    s.logs.Sub(topics...),
	}
}

// Close stops listening for events
func (l *Logger) Close() {
	go func() {
		l.ps.Unsub(l.sub, l.topics...)
		close(l.Logs)
	}()
}

// Listen listens events that match filter
func (l *Logger) Listen() {
	for v := range l.sub {
		if log, ok := v.(*OrchestratorLog); ok {
			l.Logs <- log
		}
	}
}
