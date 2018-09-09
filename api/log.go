package api

import "github.com/mesg-foundation/core/service"

// ServiceLogsFilter is a filter func for filtering service logs.
type ServiceLogsFilter func(*logLogger)

// ServiceLogsDependenciesFilter returns a dependency filter.
func ServiceLogsDependenciesFilter(dependencies ...string) ServiceLogsFilter {
	return func(s *logLogger) {
		s.dependencies = dependencies
	}
}

// ServiceLogs gives logs for all dependencies or one when specified with filters of service serviceID.
func (a *API) ServiceLogs(serviceID string, filters ...ServiceLogsFilter) ([]*service.Logs, error) {
	return newLogLogger(a, filters...).logs(serviceID)
}
