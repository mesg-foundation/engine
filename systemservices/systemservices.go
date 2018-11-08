// Package systemservices is responsible to manage all system services
// by executing their tasks, reacting on their task results and events.
package systemservices

import (
	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/service"
)

// list of system services.
// these names are also relative paths of system services in the filesystem.
const (
	resolverService = "resolver"
	workflowService = "workflow"
)

// systemServicesList is the list of system services.
// system services will be created from this list.
var systemServicesList = []string{
	resolverService,
	workflowService,
}

// systemService represents a system service.
type systemService struct {
	*service.Service

	// name is the unique name of system service.
	// it's also the relative path of system service in the filesystem.
	name string
}

// SystemServices is managing all system services.
// It is responsible to start all system services when the core start.
// It reads the services' ID from the config package.
// All system services should runs all the time.
// Any interaction with the system services are done by using the api package.
type SystemServices struct {
	api *api.API

	// absolute path of system services dir.
	systemServicesPath string

	// all deployed system services
	services []*systemService
}

// New creates a new SystemServices instance.
// It accepts an instance of the api package.
// It accepts the system services path.
// It reads the services' ID from the config package.
// It starts all system services.
// It waits for all system services to run.
// If services' ID are not in the config, it should return an error.
// If services doesn't start properly, it should return an error.
func New(api *api.API, systemServicesPath string) (*SystemServices, error) {
	s := &SystemServices{
		api:                api,
		systemServicesPath: systemServicesPath,
	}

	for _, name := range systemServicesList {
		s.services = append(s.services, &systemService{name: name})
	}

	if err := s.deployServices(); err != nil {
		return nil, err
	}
	if err := s.startServices(); err != nil {
		return nil, err
	}
	return s, nil
}

// ResolverServiceID returns resolver service id.
func (s *SystemServices) ResolverServiceID() string {
	return s.getServiceID(resolverService)
}

// WorkflowServiceID returns workflow service's id.
func (s *SystemServices) WorkflowServiceID() string {
	return s.getServiceID(workflowService)
}
