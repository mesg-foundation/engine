package api

import "github.com/mesg-foundation/core/service"

// ListServicesFilter is a filter func for listing services.
type ListServicesFilter func(*serviceLister)

// ListServices returns services matches with filters.
func (a *API) ListServices(filters ...ListServicesFilter) ([]*service.Service, error) {
	return newServiceLister(a, filters...).List()
}
