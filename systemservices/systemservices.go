package systemservices

import (
	"github.com/mesg-foundation/core/api"
)

// SystemServices is managing all system services.
// It is responsible to start all system services when the core start.
// It reads the services' ID from the config package.
// All system services should runs all the time.
// Any interaction with the system services are done by using the api package.
type SystemServices struct {
	api                *api.API
	systemServicesPath string
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
	_, err := s.deploySystemServices()
	if err != nil {
		return nil, err
	}
	return s, nil
}
