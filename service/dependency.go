package service

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/spf13/viper"
)

func extractPorts(dependency *Dependency) (ports []container.Port) {
	ports = make([]container.Port, len(dependency.Ports))
	for i, p := range dependency.Ports {
		split := strings.Split(p, ":")
		from, _ := strconv.ParseUint(split[0], 10, 64)
		to := from
		if len(split) > 1 {
			to, _ = strconv.ParseUint(split[1], 10, 64)
		}
		ports[i] = container.Port{
			Target:    uint32(to),
			Published: uint32(from),
		}
	}
	return
// DependencyFromService represents a Dependency, with a pointer to its service and its name
type DependencyFromService struct {
	*Dependency
	Service *Service
	Name    string
}

func extractVolumes(service *Service, dependency *Dependency, details dependencyDetails) (volumes []container.Mount) {
	volumes = make([]container.Mount, 0)
	for _, volume := range dependency.Volumes {
		path := filepath.Join(details.namespace, details.dependencyName, volume)
		source := filepath.Join(viper.GetString(config.ServicePathHost), path)
		volumes = append(volumes, container.Mount{
			Source: source,
			Target: volume,
// DependenciesFromService returns the an array of DependencyFromService
func (s *Service) DependenciesFromService() (d []*DependencyFromService) {
	for name, dependency := range s.GetDependencies() {
		d = append(d, &DependencyFromService{
			Dependency: dependency,
			Service:    s,
			Name:       name,
		})
		os.MkdirAll(filepath.Join(viper.GetString(config.ServicePathDocker), path), os.ModePerm)
	}
	for _, dep := range dependency.Volumesfrom {
		for _, volume := range service.Dependencies[dep].Volumes {
			path := filepath.Join(details.namespace, dep, volume)
			source := filepath.Join(viper.GetString(config.ServicePathHost), path)
			volumes = append(volumes, container.Mount{
				Source: source,
				Target: volume,
			})
			os.MkdirAll(filepath.Join(viper.GetString(config.ServicePathDocker), path), os.ModePerm)
		}
	}
	return
}

func dependencyStatus(namespace string, dependencyName string) (status StatusType) {
	dockerStatus, err := container.ServiceStatus([]string{namespace, dependencyName})
	if err != nil {
		panic(err) //TODO: that's ugly
	}
	status = STOPPED
	if dockerStatus == container.RUNNING {
		status = RUNNING
	}
	return
}
