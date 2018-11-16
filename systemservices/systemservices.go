// Package systemservices is responsible to deploy & start all
// system services and provide their service ids.
package systemservices

import (
	"fmt"

	"github.com/mesg-foundation/core/service"
)

// list of system services.
// these names are also relative paths of system services in the filesystem.
const (
	ResolverService = "resolver"
	WorkflowService = "workflow"
)

// SystemServicesList is the list of system services.
// system services will be created from this list.
var SystemServicesList = []string{
	ResolverService,
	WorkflowService,
}

// systemService represents a system service.
type systemService struct {
	service *service.Service

	// name is the unique name of system service.
	// it's also the relative path of system service in the filesystem.
	name string
}

// SystemServices is managing all system services.
// It is responsible to start all system services when the core start.
// All system services should run all the time.
// Any interaction with the system services are done by using the api package.
type SystemServices struct {
	services []*systemService
}

// New creates a new SystemServices instance.
func New() *SystemServices {
	return &SystemServices{}
}

// RegisterSystemService adds a deployed system service in the systemservices manager
func (s *SystemServices) RegisterSystemService(name string, service *service.Service) error {
	for _, ss := range s.services {
		if ss.name == name {
			return fmt.Errorf("System service already registered")
		}
	}
	s.services = append(s.services, &systemService{
		name:    name,
		service: service,
	})
	return nil
}

// GetServiceID returns the service id of a system service that matches with name.
// name compared with the unique name/relative path of system service.
func (s *SystemServices) GetServiceID(name string) (string, error) {
	for _, srv := range s.services {
		if srv.name == name {
			return srv.service.ID, nil
		}
	}
	return "", &SystemServiceNotFoundError{Name: name}
}

// ResolverServiceID returns resolver system service's id.
func (s *SystemServices) ResolverServiceID() (string, error) {
	return s.GetServiceID(ResolverService)
}

// WorkflowServiceID returns workflow service's id.
func (s *SystemServices) WorkflowServiceID() (string, error) {
	return s.GetServiceID(WorkflowService)
}
