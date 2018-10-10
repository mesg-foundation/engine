package systemservices

import (
	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/systemservices/resolver"
)

// SystemServices is managing all system services.
// It is responsible to start all system services when the core start.
// It reads the services' ID from the config package.
// All system services should runs all the time.
// Any interaction with the system services are done by using the api package.
type SystemServices struct {
	api                *api.API
	systemServicesPath string

	// system services
	resolverService *resolver.Resolver
}

// New creates a new SystemServices instance.
// It accepts an instance of the api package.
// It accepts the system services path.
// It reads the services' ID from the config package.
// It starts all system services.
// It waits for all system services to run.
// If services' ID are not in the config, it should return an error.
// IF the services don't start properly, it should return an error.
func New(api *api.API, systemServicesPath string) (*SystemServices, error) {
	s := &SystemServices{
		api:                api,
		systemServicesPath: systemServicesPath,
	}
	services, err := s.deploySystemServices()
	if err != nil {
		return nil, err
	}
	if err := s.startServices(services); err != nil {
		return nil, err
	}
	return s, s.initSystemServices(services)
}

// systemService type used to create a key, service name pairs.
type systemService int

const (
	resolverService systemService = iota
)

// systemServices keeps the system services' names.
var systemServices = map[systemService]string{
	resolverService: "System Resolver Service",
}

// initSystemServices initializes all system services.
func (s *SystemServices) initSystemServices(services []*service.Service) error {
	var err error

	// init resolver system service.
	resolverServiceID := s.getServiceID(services, systemServices[resolverService])
	if resolverServiceID == "" {
		return &systemServiceNotFound{name: systemServices[resolverService]}
	}
	s.resolverService, err = resolver.New(resolverServiceID, s.api)
	if err != nil {
		return err
	}

	return nil
}
