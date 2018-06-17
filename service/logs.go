package service

import (
	"io"

	"github.com/mesg-foundation/core/container"
)

// Logs return the service's docker service logs. Optionally only show the logs of a given dependency
func (service *Service) Logs(onlyForDependency string) (readers []io.ReadCloser, err error) {
	for _, dep := range service.DependenciesFromService() {
		if onlyForDependency == "" || onlyForDependency == "*" || onlyForDependency == dep.Name {
			var reader io.ReadCloser
			reader, err = container.ServiceLogs(dep.namespace())
			if err != nil {
				break
			}
			readers = append(readers, reader)
		}
	}
	return
}
