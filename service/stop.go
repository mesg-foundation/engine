package service

// Stop a service
func (service *Service) Stop() (err error) {
	if service.IsStopped() {
		return
	}
	for name, dependency := range service.Dependencies {
		err = dependency.Stop(service.namespace(), name)
		if err != nil {
			break
		}
	}
	return
}
