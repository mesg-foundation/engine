package applicationtest

type ServiceStart struct {
	serviceID string
}

func (s *ServiceStart) ServiceID() string {
	return s.serviceID
}
