package sdk

import (
	"crypto/sha1"
	"encoding/hex"
	"io"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/mesg-foundation/core/hash"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/x/xstrings"
	"github.com/sirupsen/logrus"
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
	var (
		logs       []*service.Log
		isNoFilter = len(f.dependencies) == 0
	)
	for _, d := range append(s.Dependencies, s.Configuration) {
		// Service.Configuration can be nil so, here is a check for it.
		if d == nil {
			continue
		}
		if isNoFilter || xstrings.SliceContains(f.dependencies, d.Key) {
			var r io.ReadCloser
			r, err := sdk.container.ServiceLogs(dependencyNamespace(serviceNamespace(s.Hash), d.Key))
			if err != nil {
				return nil, err
			}
			rstd, sw := io.Pipe()
			rerr, ew := io.Pipe()
			go func(dstout, dsterr io.Writer, r io.ReadCloser) {
				if _, err := stdcopy.StdCopy(dstout, dsterr, r); err != nil {
					r.Close()
					logrus.Errorln(err)
				}
			}(sw, ew, r)
			logs = append(logs, &service.Log{
				Dependency: d.Key,
				Standard:   rstd,
				Error:      rerr,
			})
		}
	}
	return logs, nil
}

// serviceNamespace returns the namespace of the service.
func serviceNamespace(hash hash.Hash) []string {
	sum := sha1.Sum(hash)
	return []string{hex.EncodeToString(sum[:])}
}

// dependencyNamespace builds the namespace of a dependency.
func dependencyNamespace(serviceNamespace []string, dependencyKey string) []string {
	return append(serviceNamespace, dependencyKey)
}
