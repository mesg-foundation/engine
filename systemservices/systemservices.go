// Package systemservices is responsible to deploy & start all
// system services and provide their service ids.
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
// All system services should run all the time.
// Any interaction with the system services are done by using the api package.
type SystemServices struct {
	d *deployer.Deployer
}

// New creates a new SystemServices instance.
func New(d *deployer.Deployer) (*SystemServices, error) {
	s := &SystemServices{d: d}
	return s, s.d.Deploy(systemServicesList)
}

// ResolverServiceID returns resolver system service's id.
func (s *SystemServices) ResolverServiceID() string {
	return s.d.GetServiceID(ResolverService)
}
