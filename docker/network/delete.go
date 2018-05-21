package network

import (
	"github.com/mesg-foundation/core/docker"
)

// Delete deletes a Docker Network associated with a namespace
func Delete(namespace string) (err error) {
	client, err := docker.Client()
	if err != nil {
		return
	}
	network, err := Find(namespace)
	if err != nil {
		return
	}
	return client.RemoveNetwork(network.ID)
}
