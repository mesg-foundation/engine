package service

import (
	"sync"

	"github.com/mesg-foundation/core/container"
)

// Stop stops a service.
func (s *Service) Stop() error {
	status, err := s.Status()
	if err != nil || status == STOPPED {
		return err
	}

	if err := s.StopDependencies(); err != nil {
		return err
	}
	return s.container.DeleteNetwork(s.namespace(), container.EventDestroy)
}

// StopDependencies stops all dependencies.
func (s *Service) StopDependencies() error {
	var mutex sync.Mutex
	var wg sync.WaitGroup
	var err error
	for _, dep := range s.Dependencies {
		wg.Add(1)
		go func(d *Dependency) {
			defer wg.Done()
			errStop := d.Stop()
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
func (d *Dependency) Stop() error {
	status, err := d.Status()
	if err != nil || status == container.STOPPED {
		return err
	}
	return d.service.container.StopService(d.namespace())
}
