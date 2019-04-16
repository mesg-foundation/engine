package service

import (
	"io"

	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/x/xstrings"
)

// Log holds log streams of dependency.
type Log struct {
	Dependency string
	Standard   io.ReadCloser
	Error      io.ReadCloser
}

// Logs gives service's logs and applies dependencies filter to filter logs.
// if dependencies has a length of zero all dependency logs will be provided.
func (s *Service) Logs(c container.Container, dependencies ...string) ([]*Log, error) {
	var (
		logs       []*Log
		isNoFilter = len(dependencies) == 0
	)
	addLog := func(dep *Dependency, name string) error {
		if isNoFilter || xstrings.SliceContains(dependencies, name) {
			rstd, rerr, err := dep.Logs(c, s.namespace())
			if err != nil {
				return err
			}
			logs = append(logs, &Log{
				Dependency: name,
				Standard:   rstd,
				Error:      rerr,
			})
		}
		return nil
	}
	if s.Configuration != nil {
		if err := addLog(s.Configuration, MainServiceKey); err != nil {
			return nil, err
		}
	}
	for _, dep := range s.Dependencies {
		if err := addLog(dep, dep.Key); err != nil {
			return nil, err
		}
	}
	return logs, nil
}
