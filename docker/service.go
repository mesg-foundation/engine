package docker

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types/swarm"
	godocker "github.com/fsouza/go-dockerclient"
)

// Service returns the Docker Service
func Service(namespace string, name string) (dockerService swarm.Service, err error) {
	ctx := context.Background()
	client, err := Client()
	if err != nil {
		return
	}
	dockerServices, err := client.ListServices(godocker.ListServicesOptions{
		Filters: map[string][]string{
			"name": []string{strings.Join([]string{namespace, name}, "_")},
		},
		Context: ctx,
	})
	if err != nil {
		return
	}
	dockerService = serviceMatch(dockerServices, namespace, name)
	return
}

func serviceMatch(dockerServices []swarm.Service, namespace string, name string) (dockerService swarm.Service) {
	for _, service := range dockerServices {
		if service.Spec.Annotations.Name == strings.Join([]string{namespace, name}, "_") {
			dockerService = service
			break
		}
	}
	return
}
