package systemservices

import "github.com/mesg-foundation/core/systemservices/resolver"

// Resolver returns the Resolver instance using the running Resolver service.
func (s *SystemServices) Resolver() *resolver.Resolver {
	return s.resolverService
}
