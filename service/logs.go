package service

import (
	"io"

	"github.com/mesg-foundation/core/container"
)

// Logs return the service's docker service logs. Optionally only show the logs of a given dependency
func (service *Service) Logs(onlyForDependency string) (readers []io.ReadCloser, err error) {
	for depName := range service.GetDependencies() {
		if onlyForDependency == "" || onlyForDependency == "*" || onlyForDependency == depName {
			namespace := []string{service.namespace(), depName} //TODO: refacto namespace
			var reader io.ReadCloser
			reader, err = container.ServiceLogs(namespace)
			if err != nil {
				break
			}
			readers = append(readers, reader)
		}
	}
	return
}
