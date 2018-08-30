package service

import (
	"github.com/mesg-foundation/core/container"
)

// StatusType of the service.
type StatusType uint

// Possible statuses for service.
const (
	UNKNOWN StatusType = iota
	STOPPED
	STARTING
	PARTIAL
	RUNNING
)

var containerStatusTypeMappings = map[container.StatusType]StatusType{
	container.UNKNOWN:  UNKNOWN,
	container.STOPPED:  STOPPED,
	container.STARTING: STARTING,
	container.RUNNING:  RUNNING,
}

// Status returns StatusType of all dependency.
func (service *Service) Status() (StatusType, error) {
	var (
		dependencies = service.DependenciesFromService()
		statuses     map[container.StatusType]bool
	)
	for _, d := range dependencies {
		status, err := d.Status()
		if err != nil {
			return UNKNOWN, err
		}
		statuses[status] = true
	}

	switch len(statuses) {
	case 0:
		return STOPPED, nil
	case 1:
		for status := range statuses {
			return containerStatusTypeMappings[status], nil
		}
	default:
		return PARTIAL, nil
	}
	panic("not reached")
}

// Status returns StatusType of dependency's container.
func (dependency *DependencyFromService) Status() (container.StatusType, error) {
	return defaultContainer.Status(dependency.namespace())
}

// ListRunning returns all the running services.2
// TODO: should move to another file
func ListRunning() ([]string, error) {
	services, err := defaultContainer.ListServices("mesg.hash")
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
