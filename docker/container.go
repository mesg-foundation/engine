package docker

import (
	"context"

	godocker "github.com/fsouza/go-dockerclient"
)

// FindContainer returns a running docker container if exist
func FindContainer(name string) (*godocker.APIContainers, error) {
	client, err := Client()
	if err != nil {
		return nil, nil
	}
	res, err := client.ListContainers(godocker.ListContainersOptions{
		Context: context.Background(),
		Limit:   1,
		Filters: map[string][]string{
			"ancestor": []string{name},
			"status":   []string{"running"},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return &res[0], nil
}
