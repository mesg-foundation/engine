package systemservices

import (
	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/systemservices/resolver"
)

// SystemServices is managing all system services.
// It is responsible to start all system services when the core start.
// It reads the services' ID from the config package.
// All system services should runs all the time.
// Any interaction with the system services are done by using the API package.
type SystemServices struct{}

// New creates a new SystemServices instance.
// It accepts an instance of the API package.
// It reads the services' ID from the config package.
// It starts all system services.
// It waits for all system services to run.
// If services' ID are not in the config, it should return an error.
// IF the services don't start properly, it should return an error.
func New(api *api.API) (*SystemServices, error) {
	return nil, nil
}

// Resolver returns the Resolver instance using the running Resolver service.
func (ss *SystemServices) Resolver() (*resolver.Resolver, error) {
	return nil, nil
}
