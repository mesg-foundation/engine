package service

import (
	"github.com/mesg-foundation/core/container"
)

// StatusType of the service.
type StatusType uint

// Possible statuses for service.
const (
	STOPPED StatusType = 0
	PARTIAL StatusType = 1
	RUNNING StatusType = 2
)

// Status returns StatusType of all dependency.
func (service *Service) Status() (StatusType, error) {
	status := STOPPED
	allRunning := true
	for _, dependency := range service.DependenciesFromService() {
		depStatus, err := dependency.Status()
		if err != nil {
			return status, err
		}
		if depStatus == container.RUNNING {
			status = RUNNING
		} else {
			allRunning = false
		}
	}
	if status == RUNNING && !allRunning {
		status = PARTIAL
	}
	return status, nil
}

// Status returns StatusType of dependency's container.
func (dependency *DependencyFromService) Status() (container.StatusType, error) {
	return container.ServiceStatus(dependency.namespace())
}

// ListRunning returns all the running services.2
// TODO: should move to another file
func ListRunning() ([]string, error) {
	services, err := container.ListServices("mesg.hash")
	if err != nil {
		return nil, err
	}
	mapRes := make(map[string]uint)
	for _, service := range services {
		serviceName := service.Spec.Annotations.Labels["mesg.hash"]
		mapRes[serviceName]++
	}
	res := make([]string, 0, len(mapRes))
	for k := range mapRes {
		res = append(res, k)
	}
	return res, nil
}
