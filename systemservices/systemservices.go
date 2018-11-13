// Package systemservices is responsible to manage all system services
// by executing their tasks, reacting on their task results and events.
package systemservices

import "github.com/mesg-foundation/core/systemservices/deployer"

// list of system services.
// these names are also relative paths of system services in the filesystem.
const (
	ResolverService = "resolver"
)

// systemServicesList is the list of system services.
// system services will be created from this list.
var systemServicesList = []string{
	ResolverService,
}

// SystemServices is managing all system services.
// It is responsible to start all system services when the core start.
// It reads the services' ID from the config package.
// All system services should runs all the time.
// Any interaction with the system services are done by using the api package.
type SystemServices struct {
	d *deployer.Deployer
}

// New creates a new SystemServices instance.
func New(d *deployer.Deployer) (*SystemServices, error) {
	s := &SystemServices{d: d}
	return s, s.d.Deploy(systemServicesList)
}

// ResolverServiceID returns resolver service id.
func (s *SystemServices) ResolverServiceID() string {
	return s.d.GetServiceID(ResolverService)
}
