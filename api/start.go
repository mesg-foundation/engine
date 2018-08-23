package api

// StartService starts service serviceID.
func (a *API) StartService(serviceID string) error {
	return newServiceStarter(a).Start(serviceID)
}
