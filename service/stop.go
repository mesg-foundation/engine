package service

import (
	"github.com/mesg-foundation/core/container"
)

// Stop a service
func (service *Service) Stop() (err error) {
	if service.IsStopped() {
		return
	}
	for name, dependency := range service.GetDependencies() {
		err = dependency.Stop(service.namespace(), name)
		if err != nil {
			break
		}
	}
	if err == nil { // didnt exit the loop
		err = container.DeleteNetwork([]string{service.namespace()})
	}
	return
}

// Stop a dependency
func (dependency *Dependency) Stop(namespace string, dependencyName string) (err error) {
	if !dependency.IsRunning(namespace, dependencyName) {
		return
	}
	err = container.StopService([]string{namespace, dependencyName})
	return
}
