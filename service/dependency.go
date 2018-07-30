package service

import (
	"sort"
)

// DependencyFromService represents a Dependency with a pointer to its service and its name.
type DependencyFromService struct {
	*Dependency
	Service *Service
	Name    string
}

// DependenciesFromService returns a slice of DependencyFromService.
func (s *Service) DependenciesFromService() []*DependencyFromService {
	var keys []string
	for key := range s.Dependencies {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	d := make([]*DependencyFromService, len(keys))
	for _, key := range keys {
		dependency := s.Dependencies[key]
		d = append(d, &DependencyFromService{
			Dependency: dependency,
			Service:    s,
			Name:       key,
		})
	}
	return d
}
