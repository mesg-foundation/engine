package api

import (
	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/service"
)

// StartService starts service serviceID.
func (a *API) StartService(serviceID string) error {
	return newServiceStarter(a).Start(serviceID)
}

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
	sr, err := services.Get(serviceID)
	if err != nil {
		return err
	}
	sr, err = service.FromService(sr, service.ContainerOption(s.api.container))
	if err != nil {
		return err
	}
	_, err = sr.Start()
	return err
}
