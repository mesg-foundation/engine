package api

import (
	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/service"
)

// logLogger provides functionalities to get service logs.
type logLogger struct {
	// dependencies used to get only logs from requested dependencies.
	dependencies []string

	api *API
}

// newLogLogger creates a new logLogger with given api and dependency filters.
func newLogLogger(api *API, filters ...ServiceLogsFilter) *logLogger {
	l := &logLogger{
		api: api,
	}
	for _, filter := range filters {
		filter(l)
	}
	return l
}

// logs gives logs of service serviceID and applies dependency filters to filter logs.
func (l *logLogger) logs(serviceID string) ([]*service.Log, error) {
	s, err := services.Get(serviceID)
	if err != nil {
		return nil, err
	}
	return s.Logs(l.dependencies...)
}
