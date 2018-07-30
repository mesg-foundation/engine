package service

import (
	"sync"

	"github.com/mesg-foundation/core/container"
)

// Stop stops a service.
func (service *Service) Stop() error {
	status, err := service.Status()
	if err != nil || status == STOPPED {
		return err
	}

	if err := service.StopDependencies(); err != nil {
		return err
	}
	return container.DeleteNetwork(service.namespace())
}

// StopDependencies stops all dependencies.
func (service *Service) StopDependencies() error {
	var mutex sync.Mutex
	var wg sync.WaitGroup
	var err error
	for _, dependency := range service.DependenciesFromService() {
		wg.Add(1)
		go func(d *DependencyFromService) {
			defer wg.Done()
			errStop := d.Stop()
			mutex.Lock()
			defer mutex.Unlock()
			if errStop != nil && err == nil {
				err = errStop
			}
		}(dependency)
	}
	wg.Wait()
	return err
}

// Stop stops a dependency.
func (dependency *DependencyFromService) Stop() error {
	status, err := dependency.Status()
	if err != nil || status == container.STOPPED {
		return err
	}
	return container.StopService(dependency.namespace())
}
