package service

import (
	"sync"

	"github.com/mesg-foundation/core/container"
)

// Stop a service
func (service *Service) Stop() (err error) {
	if service.IsStopped() {
		return
	}
	err = service.StopDependencies()
	if err != nil {
		return
	}
	err = container.DeleteNetwork([]string{service.namespace()})
	if err != nil {
		return
	}
	return
}

// StopDependencies stops all dependencies
func (service *Service) StopDependencies() (err error) {
	var mutex sync.Mutex
	var wg sync.WaitGroup
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
	return
}

// Stop a dependency
func (dependency *DependencyFromService) Stop() (err error) {
	if !dependency.IsRunning() {
		return
	}
	err = container.StopService(dependency.namespace())
	if err != nil {
		return
	}
	err = container.WaitForContainerStatus(dependency.namespace(), container.STOPPED)
	return
}
