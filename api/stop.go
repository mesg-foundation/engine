package api

// StopService stops service serviceID.
func (a *API) StopService(serviceID string) error {
	return newServiceStopper(a).Stop(serviceID)
}
