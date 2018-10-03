package api

import (
	"github.com/mesg-foundation/core/service"
)

// GetService returns service serviceID.
func (a *API) GetService(serviceID string) (*service.Service, error) {
	return newServiceGetter(a).Get(serviceID)
}

// serviceGetter provides functionalities to get a MESG service.
type serviceGetter struct {
	api *API
}

// newServiceGetter creates a new serviceGetter with given api.
func newServiceGetter(api *API) *serviceGetter {
	return &serviceGetter{
		api: api,
	}
}

// Get returns service serviceID.
func (g *serviceGetter) Get(serviceID string) (*service.Service, error) {
	s, err := g.api.db.Get(serviceID)
	if err != nil {
		return nil, err
	}
	return service.FromService(s, service.ContainerOption(g.api.container))
}
