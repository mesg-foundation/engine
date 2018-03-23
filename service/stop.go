package service

import (
	"context"

	docker "github.com/fsouza/go-dockerclient"
)

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
	if err == nil { // didnt exit the loop
		err = deleteNetwork(service.namespace())
	}
	return
}

// Stop a dependency
func (dependency Dependency) Stop(namespace string, dependencyName string) (err error) {
	ctx := context.Background()
	cli, err := dockerCli()
	if err != nil {
		return
	}
	if !dependency.IsRunning(namespace, dependencyName) {
		return
	}
	dockerService, err := getDockerService(namespace, dependencyName)
	if err == nil && dockerService.ID != "" {
		err = cli.RemoveService(docker.RemoveServiceOptions{
			ID:      dockerService.ID,
			Context: ctx,
		})
	}
	return
}
