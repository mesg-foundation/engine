package keeper

import (
	"github.com/go-kit/kit/metrics"
	prometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

// M is global metric variable.
var M = NewMetric()

// Metric is struct to keep counters for execution statuses.
type Metric struct {
	Created    metrics.Counter
	InProgress metrics.Counter
	Updated    metrics.Counter
	Completed  metrics.Counter
}

// NewMetric creates a counters metric.
func NewMetric() *Metric {
	return &Metric{
		Created: prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "mesg",
			Subsystem: "execution",
			Name:      "created",
			Help:      "executions created",
		}, []string{}),
		InProgress: prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "mesg",
			Subsystem: "execution",
			Name:      "in_progress",
			Help:      "executions in progress",
		}, []string{}),
		Updated: prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "mesg",
			Subsystem: "execution",
			Name:      "updated",
			Help:      "executions updated",
		}, []string{}),
		Completed: prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "mesg",
			Subsystem: "execution",
			Name:      "completed",
			Help:      "executions completed",
		}, []string{}),
	}
}
