package api

import (
	"github.com/mesg-foundation/core/database/services"
)

// serviceStarter provides functionalities to start a MESG service.
type serviceStarter struct {
	api *API
}

// newServiceStarter creates a new serviceStarter with given api.
func newServiceStarter(api *API) *serviceStarter {
	return &serviceStarter{
		api: api,
	}
}

// Start starts service serviceID.
func (s *serviceStarter) Start(serviceID string) error {
	service, err := services.Get(serviceID)
	if err != nil {
		return err
	}
	_, err = service.Start()
	return err
}
