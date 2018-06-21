package service

import (
	"sort"
)

// DependencyFromService represents a Dependency, with a pointer to its service and its name
type DependencyFromService struct {
	*Dependency
	Service *Service
	Name    string
}

// DependenciesFromService returns the an array of DependencyFromService
func (s *Service) DependenciesFromService() (d []*DependencyFromService) {
	var keys []string
	for key := range s.Dependencies {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		dependency := s.Dependencies[key]
		d = append(d, &DependencyFromService{
			Dependency: dependency,
			Service:    s,
			Name:       key,
		})
	}
	return
}
