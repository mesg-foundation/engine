package service

// Start a service
func (service *Service) Start() (err error) {
	for name, dependency := range service.Dependencies {
		err = dependency.Start(name, service.namespace())
	}
	// Disgrasfully close the service because there is an error
	if err != nil {
		service.Stop()
	}
	return
}
