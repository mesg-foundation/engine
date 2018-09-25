package api

import (
	"github.com/mesg-foundation/core/service"
)

// ServiceLogsFilter is a filter func for filtering service logs.
type ServiceLogsFilter func(*logLogger)

// ServiceLogsDependenciesFilter returns a dependency filter.
func ServiceLogsDependenciesFilter(dependencies ...string) ServiceLogsFilter {
	return func(s *logLogger) {
		s.dependencies = dependencies
	}
}

// ServiceLogs gives logs for all dependencies or one when specified with filters of service serviceID.
func (a *API) ServiceLogs(serviceID string, filters ...ServiceLogsFilter) ([]*service.Log, error) {
	return newLogLogger(a, filters...).logs(serviceID)
}

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
	s, err := l.api.db.Get(serviceID)
	if err != nil {
		return nil, err
	}
	s, err = service.FromService(s, service.ContainerOption(l.api.container))
	if err != nil {
		return nil, err
	}
	return s.Logs(l.dependencies...)
}
