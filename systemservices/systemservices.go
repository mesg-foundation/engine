// Package systemservices is responsible to manage all system services
// by executing their tasks, reacting on their task results and events.
package systemservices

import (
	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/systemservices/resolver"
)

// list of system services.
// these names are also relative paths of system services in the filesystem.
const (
	resolverService = "resolver"
)

// systemServicesList is the list of system services.
// system services will be created from this list.
var systemServicesList = []string{
	resolverService,
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

	// system services.
	resolverService *resolver.Resolver
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

	var services []*systemService
	for _, name := range systemServicesList {
		services = append(services, &systemService{name: name})
	}

	if err := s.deployServices(services); err != nil {
		return nil, err
	}
	if err := s.startServices(services); err != nil {
		return nil, err
	}
	return s, s.initServices(services)
}

// Resolver returns the Resolver instance using the running Resolver service.
func (s *SystemServices) Resolver() *resolver.Resolver {
	return s.resolverService
}

// initServices initializes all system services.
func (s *SystemServices) initServices(services []*systemService) error {
	// init resolver system service.
	s.resolverService = resolver.New(s.getServiceID(services, resolverService), s.api)
	return nil
}
