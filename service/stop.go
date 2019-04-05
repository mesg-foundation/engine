package service

import (
	"sync"

	"github.com/mesg-foundation/core/container"
)

// Stop stops a service.
func (s *Service) Stop(c container.Container) error {
	status, err := s.Status(c)
	if err != nil || status == STOPPED {
		return err
	}

	if err := s.StopDependencies(c); err != nil {
		return err
	}
	return c.DeleteNetwork(s.namespace())
}

// StopDependencies stops all dependencies.
func (s *Service) StopDependencies(c container.Container) error {
	var mutex sync.Mutex
	var wg sync.WaitGroup
	var err error
	for _, dep := range s.Dependencies {
		wg.Add(1)
		go func(d *Dependency) {
			defer wg.Done()
			errStop := d.Stop(c)
			mutex.Lock()
			defer mutex.Unlock()
			if errStop != nil && err == nil {
				err = errStop
			}
		}(dep)
	}
	wg.Wait()
	return err
}

// Stop stops a dependency.
func (d *Dependency) Stop(c container.Container) error {
	return c.StopService(d.namespace())
}
