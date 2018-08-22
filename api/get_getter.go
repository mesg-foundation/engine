package api

import (
	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/service"
)

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
	service, err := services.Get(serviceID)
	return &service, err
}
