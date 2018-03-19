package service

// Stop a service
func (service *Service) Stop() (err error) {
	if service.IsStopped() {
		return
	}
	for name, dependency := range service.Dependencies {
		err = dependency.Stop(name, service.namespace())
		if err != nil {
			break
		}
	}
	return
}
