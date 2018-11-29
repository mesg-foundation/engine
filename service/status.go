package service

import (
	"github.com/mesg-foundation/core/config"
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

func (s StatusType) String() string {
	switch s {
	case STOPPED:
		return "STOPPED"
	case STARTING:
		return "STARTING"
	case PARTIAL:
		return "PARTIAL"
	case RUNNING:
		return "RUNNING"
	default:
		return "UNKNOWN"
	}
}

var containerStatusTypeMappings = map[container.StatusType]StatusType{
	container.UNKNOWN:  UNKNOWN,
	container.STOPPED:  STOPPED,
	container.STARTING: STARTING,
	container.RUNNING:  RUNNING,
}

// Status returns StatusType of all dependency.
func (s *Service) Status() (StatusType, error) {
	statuses := make(map[container.StatusType]bool)
	for _, dep := range s.Dependencies {
		status, err := dep.Status()
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
func (d *Dependency) Status() (container.StatusType, error) {
	return d.service.container.Status(d.namespace())
}

// ListRunning returns all the running services.2
// TODO: should move to another file
func ListRunning() ([]string, error) {
	cfg, err := config.Global()
	if err != nil {
		return nil, err
	}

	// TODO(ilgooz): remove this line after ListRunning refactored.
	c, err := container.New()
	if err != nil {
		return nil, err
	}
	services, err := c.ListServices("mesg.hash", "mesg.core="+cfg.Core.Name)
	if err != nil {
		return nil, err
	}
	// Make service list unique. One mesg service can have multiple docker service.
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
