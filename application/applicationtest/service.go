package applicationtest

// ServiceStart holds information about a service start request.
type ServiceStart struct {
	serviceID string
}

// ServiceID returns the started service's id.
func (s *ServiceStart) ServiceID() string {
	return s.serviceID
}
