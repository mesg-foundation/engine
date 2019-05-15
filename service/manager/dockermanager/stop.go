package dockermanager

import (
	"sync"

	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/x/xerrors"
)

// Stop stops a service.
func (m *DockerManager) Stop(s *service.Service) error {
	status, err := m.Status(s)
	if err != nil || status == service.STOPPED {
		return err
	}

	var (
		wg   sync.WaitGroup
		errs xerrors.SyncErrors
	)
	for _, d := range append([]*service.Dependency{s.Configuration}, s.Dependencies...) {
		// Service.Configuration can be nil so, here is a check for it.
		if d == nil {
			continue
		}
		wg.Add(1)
		go func(namespace []string) {
			defer wg.Done()
			if err := m.c.StopService(namespace); err != nil {
				errs.Append(err)
			}
		}(d.Namespace(s.Namespace()))
	}
	wg.Wait()
	if err := errs.ErrorOrNil(); err != nil {
		return err
	}

	return m.c.DeleteNetwork(s.Namespace())
}
