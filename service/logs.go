package service

import (
	"io"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/x/xstrings"
	"github.com/sirupsen/logrus"
)

// Log holds log streams of dependency.
type Log struct {
	Dependency string
	Standard   io.ReadCloser
	Error      io.ReadCloser
}

// Logs gives service's logs streams. when dependencies filter is not provided, it'll give
// logs for all dependencies otherwise it'll only give logs for specified dependencies.
// note that, service itself is also a dependency defined with special "service" key.
// in order to get service's own logs, "service" key must be included to dependencies filter.
func (s *Service) Logs(c container.Container, dependencies ...string) ([]*Log, error) {
	var (
		logs       []*Log
		isNoFilter = len(dependencies) == 0
	)
	for _, d := range append(s.Dependencies, s.Configuration) {
		// Service.Configuration can be nil so, here is a check for it.
		if d == nil {
			continue
		}
		if isNoFilter || xstrings.SliceContains(dependencies, d.Key) {
			var r io.ReadCloser
			r, err := c.ServiceLogs(d.namespace(s.namespace()))
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
			if err != nil {
				return nil, err
			}
			logs = append(logs, &Log{
				Dependency: d.Key,
				Standard:   rstd,
				Error:      rerr,
			})
		}
	}
	return logs, nil
}
