package dockermanager

import (
	"io"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/x/xstrings"
	"github.com/sirupsen/logrus"
)

// Logs gives service's logs streams. when dependencies filter is not provided, it'll give
// logs for all dependencies otherwise it'll only give logs for specified dependencies.
// note that, service itself is also a dependency defined with special "service" key.
// in order to get service's own logs, "service" key must be included to dependencies filter.
func (m *DockerManager) Logs(s *service.Service, dependencies ...string) ([]*service.Log, error) {
	var (
		logs       []*service.Log
		isNoFilter = len(dependencies) == 0
	)
	for _, d := range append(s.Dependencies, s.Configuration) {
		// Service.Configuration can be nil so, here is a check for it.
		if d == nil {
			continue
		}
		if isNoFilter || xstrings.SliceContains(dependencies, d.Key) {
			var r io.ReadCloser
			r, err := m.c.ServiceLogs(d.Namespace(s.Namespace()))
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
