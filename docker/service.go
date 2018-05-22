package docker

import (
	"context"

	"github.com/docker/docker/api/types/swarm"
	godocker "github.com/fsouza/go-dockerclient"
)

//  FindService returns the Docker Service
func FindService(name []string) (dockerService *swarm.Service, err error) {
	ctx := context.Background()
	client, err := Client()
	if err != nil {
		return
	}
	dockerServices, err := client.ListServices(godocker.ListServicesOptions{
		Filters: map[string][]string{
			"name": []string{Namespace(name)},
		},
		Context: ctx,
	})
	if err != nil {
		return
	}
	dockerService = serviceMatch(dockerServices, name)
	return
}

func serviceMatch(dockerServices []swarm.Service, name []string) (dockerService *swarm.Service) {
	for _, service := range dockerServices {
		if service.Spec.Annotations.Name == Namespace(name) {
			dockerService = &service
			break
		}
	}
	return
}
