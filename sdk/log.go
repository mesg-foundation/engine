package sdk

import (
	"github.com/mesg-foundation/core/hash"
	"github.com/mesg-foundation/core/service"
)

// serviceLogFilters keeps log filters.
type serviceLogFilters struct {
	// dependencies list is used to provide logs for only requested dependencies.
	dependencies []string
}

// ServiceLogsFilter is a filter func for filtering service logs.
type ServiceLogsFilter func(*serviceLogFilters)

// ServiceLogsDependenciesFilter returns a dependency filter.
func ServiceLogsDependenciesFilter(dependencies ...string) ServiceLogsFilter {
	return func(s *serviceLogFilters) {
		s.dependencies = dependencies
	}
}

// ServiceLogs provides logs for dependencies of service serviceID that matches with filters.
// when no dependency filters are set, all the dependencies' logs will be provided.
func (sdk *SDK) ServiceLogs(serviceHash hash.Hash, filters ...ServiceLogsFilter) ([]*service.Log, error) {
	f := &serviceLogFilters{}
	for _, filter := range filters {
		filter(f)
	}
	s, err := sdk.db.Get(serviceHash)
	if err != nil {
		return nil, err
	}
	return sdk.m.Logs(s, f.dependencies...)
}
