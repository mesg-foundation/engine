package orchestrator

import (
	"github.com/cskr/pubsub"
)

const (
	// anyActionTopic is the default topic where all logs are published to.
	anyActionTopic = "*"
)

// LogsListener provides functionalities to listen to the orchestrator's logs
type LogsListener struct {
	Logs chan *OrchestratorLog

	ps     *pubsub.PubSub
	topics []string
	sub    chan interface{}
}

// NewLogsListener creates a new PubSub that will contain action for the given topics
func (s *Orchestrator) NewLogsListener(topics ...string) *LogsListener {
	if len(topics) == 0 {
		topics = []string{anyActionTopic}
	}
	return &LogsListener{
		Logs:   make(chan *OrchestratorLog),
		ps:     s.logsPubSub,
		topics: topics,
		sub:    s.logsPubSub.Sub(topics...),
	}
}

// Close stops listening for events
func (l *LogsListener) Close() {
	go func() {
		l.ps.Unsub(l.sub, l.topics...)
		close(l.Logs)
	}()
}

// Listen listens events that match filter
func (l *LogsListener) Listen() {
	for v := range l.sub {
		if log, ok := v.(*OrchestratorLog); ok {
			l.Logs <- log
		}
	}
}
