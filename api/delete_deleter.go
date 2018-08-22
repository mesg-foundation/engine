package api

import "github.com/mesg-foundation/core/database/services"

type serviceDeleter struct {
	api *API
}

func newServiceDeleter(api *API) *serviceDeleter {
	return &serviceDeleter{
		api: api,
	}
}

// Delete stops and deletes service serviceID.
func (d *serviceDeleter) Delete(serviceID string) error {
	service, err := services.Get(serviceID)
	if err != nil {
		return err
	}
	if err := service.Stop(); err != nil {
		return err
	}
	return services.Delete(serviceID)
}
