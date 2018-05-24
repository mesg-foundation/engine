package docker

import (
	"errors"
	"time"

	godocker "github.com/fsouza/go-dockerclient"
)

const networkNamespacePrefix string = "network"

func networkNamespace(namespace []string) string {
	return Namespace(append([]string{networkNamespacePrefix}, namespace...))
}

// CreateNetwork creates a Docker Network with a namespace
func CreateNetwork(namespace []string) (network *godocker.Network, err error) {
	network, err = FindNetwork(namespace)
	if network != nil || err != nil {
		return
	}
	namespaceFlat := networkNamespace(namespace)
	client, err := Client()
	if err != nil {
		return
	}
	network, err = client.CreateNetwork(godocker.CreateNetworkOptions{
		Name:           namespaceFlat,
		CheckDuplicate: true, // Cannot have 2 network with the same name
		Driver:         "overlay",
		Labels: map[string]string{
			"com.docker.stack.namespace": namespaceFlat,
		},
	})
	return
}

// DeleteNetwork deletes a Docker Network associated with a namespace
func DeleteNetwork(namespace []string) (err error) {
	network, err := FindNetwork(namespace)
	if network == nil || err != nil {
		return
	}
	client, err := Client()
	if err != nil {
		return
	}
	return client.RemoveNetwork(network.ID)
}

// FindNetwork finds a Docker Network by a namespace
func FindNetwork(namespace []string) (network *godocker.Network, err error) {
	client, err := Client()
	if err != nil {
		return
	}
	network, err = client.NetworkInfo(networkNamespace(namespace))
	if err != nil {
		switch err.(type) {
		case *godocker.NoSuchNetwork:
			err = nil
		default:
		}
	}
	return
}

// WaitNetworkDeletion wait a network to be delete
func WaitNetworkDeletion(namespace []string, timeout time.Duration) (wait chan error) {
	start := time.Now()
	wait = make(chan error, 1)
	go func() {
		for {
			network, err := FindNetwork(namespace)
			if err != nil {
				wait <- err
				return
			}
			if network == nil {
				close(wait)
				return
			}
			diff := time.Now().Sub(start)
			if diff.Nanoseconds() >= int64(timeout) {
				wait <- errors.New("Wait too long for the network to get removed, timeout reached")
				return
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()
	return
}
