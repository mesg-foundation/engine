package docker

import (
	"context"

	godocker "github.com/fsouza/go-dockerclient"
)

// TODO: to update. make it more useful
// List all the running docker services
func List() (res []string, err error) {
	client, err := Client()
	services, err := client.ListServices(godocker.ListServicesOptions{
		Context: context.Background(),
	})
	mapRes := make(map[string]uint)
	for _, service := range services {
		serviceName := service.Spec.Annotations.Labels["mesg.service"]
		mapRes[serviceName]++
	}
	res = make([]string, 0, len(mapRes))
	for k := range mapRes {
		res = append(res, k)
	}
	return
}
