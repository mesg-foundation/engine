package service

import (
	"sync"

	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/x/xerrors"
)

// Stop stops a service.
func (s *Service) Stop(c container.Container) error {
	status, err := s.Status(c)
	if err != nil || status == STOPPED {
		return err
	}

	var (
		wg   sync.WaitGroup
		errs xerrors.SyncErrors
	)
	for _, d := range append([]*Dependency{s.Configuration}, s.Dependencies...) {
		// Service.Configuration can be nil so, here is a check for it.
		if d == nil {
			continue
		}
		wg.Add(1)
		go func(namespace []string) {
			defer wg.Done()
			if err := c.StopService(namespace); err != nil {
				errs.Append(err)
			}
		}(d.namespace(s.namespace()))
	}
	wg.Wait()
	if err := errs.ErrorOrNil(); err != nil {
		return err
	}

	return c.DeleteNetwork(s.namespace())
}
