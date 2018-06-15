package service

import (
	"github.com/mesg-foundation/core/container"
)

// StatusType of the service
type StatusType uint

// status for services
const (
	STOPPED StatusType = 1
	RUNNING StatusType = 2
	PARTIAL StatusType = 3
)

// Status returns the StatusType of all dependency of this service
func (service *Service) Status() (status StatusType, err error) {
	status = STOPPED
	allRunning := true
	for _, dependency := range service.DependenciesFromService() {
		var depStatus container.StatusType
		depStatus, err = dependency.Status()
		if err != nil {
			return
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
	return
}

// Status returns the StatusType of this dependency's container
func (dependency *DependencyFromService) Status() (status container.StatusType, err error) {
	status, err = container.ServiceStatus(dependency.namespace())
	return
}

// ListRunning all the running services
// TODO: should move to another file
func ListRunning() (res []string, err error) {
	services, err := container.ListServices("mesg.hash")
	mapRes := make(map[string]uint)
	for _, service := range services {
		serviceName := service.Spec.Annotations.Labels["mesg.hash"]
		mapRes[serviceName]++
	}
	res = make([]string, 0, len(mapRes))
	for k := range mapRes {
		res = append(res, k)
	}
	return
}
