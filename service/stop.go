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
	for name, dependency := range service.GetDependencies() {
		wg.Add(1)
		go func(d *Dependency, name string, dependencyName string) {
			defer wg.Done()
			errStop := d.Stop(name, dependencyName)
			mutex.Lock()
			defer mutex.Unlock()
			if errStop != nil && err == nil {
				err = errStop
			}
		}(dependency, service.namespace(), name)
	}
	wg.Wait()
	return
}

// Stop a dependency
func (dependency *Dependency) Stop(name string, dependencyName string) (err error) {
	if !dependency.IsRunning(name, dependencyName) {
		return
	}
	namespace := []string{name, dependencyName} //TODO: refacto namespace
	err = container.StopService(namespace)
	if err != nil {
		return
	}
	err = container.WaitForContainerStatus(namespace, container.STOPPED)
	return
}
