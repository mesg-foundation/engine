package service

import (
	"io"

	"github.com/mesg-foundation/core/container"
)

// Logs returns the service's docker service logs. Optionally only shows the logs of a given dependency.
func (service *Service) Logs(onlyForDependency string) ([]io.ReadCloser, error) {
	var readers []io.ReadCloser
	for _, dep := range service.DependenciesFromService() {
		if onlyForDependency == "" || onlyForDependency == "*" || onlyForDependency == dep.Name {
			var reader io.ReadCloser
			reader, err := container.ServiceLogs(dep.namespace())
			if err != nil {
				return readers, err
			}
			readers = append(readers, reader)
		}
	}
	return readers, nil
}
