package api

// DeleteService stops and deletes service serviceID.
func (a *API) DeleteService(serviceID string) error {
	return newServiceDeleter(a).Delete(serviceID)
}
