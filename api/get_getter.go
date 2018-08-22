package api

import (
	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/service"
)

type serviceGetter struct {
	api *API
}

func newServiceGetter(api *API) *serviceGetter {
	return &serviceGetter{
		api: api,
	}
}

// Get returns service serviceID.
func (d *serviceGetter) Get(serviceID string) (*service.Service, error) {
	service, err := services.Get(serviceID)
	return &service, err
}
