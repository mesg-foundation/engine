package service

import (
	"sort"
)

// Dependency represents a Docker container and it holds instructions about
// how it should run.
type Dependency struct {
	// Image is the Docker image.
	Image string `hash:"name:1" yaml:"image"`

	// Volumes are the Docker volumes.
	Volumes []string `hash:"name:2" yaml:"volumes"`

	// VolumesFrom are the docker volumes-from from.
	VolumesFrom []string `hash:"name:3" yaml:"volumesfrom"`

	// Ports holds ports configuration for container.
	Ports []string `hash:"name:4" yaml:"ports"`

	// Command is the Docker command which will be executed when container started.
	Command string `hash:"name:5" yaml:"command"`
}

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

	d := make([]*DependencyFromService, 0, len(keys))
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
