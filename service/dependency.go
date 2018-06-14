package service

// DependencyFromService represents a Dependency, with a pointer to its service and its name
type DependencyFromService struct {
	*Dependency
	Service *Service
	Name    string
}

// DependenciesFromService returns the an array of DependencyFromService
func (s *Service) DependenciesFromService() (d []*DependencyFromService) {
	for name, dependency := range s.GetDependencies() {
		d = append(d, &DependencyFromService{
			Dependency: dependency,
			Service:    s,
			Name:       name,
		})
	}
	return
}
