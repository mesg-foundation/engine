package docker

import (
	"errors"
	"time"

	godocker "github.com/fsouza/go-dockerclient"
)

// CreateNetwork creates a Docker Network with a namespace
func CreateNetwork(namespace []string, driver string) (network *godocker.Network, err error) {
	network, err = FindNetwork(namespace)
	if network != nil || err != nil {
		return
	}
	namespaceFlat := Namespace(namespace)
	client, err := Client()
	if err != nil {
		return
	}
	network, err = client.CreateNetwork(godocker.CreateNetworkOptions{
		Name:           namespaceFlat,
		CheckDuplicate: true, // Cannot have 2 network with the same name
		Driver:         driver,
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

// FindNetworkStrict finds a Docker Network by a namespace. If no network if found, an error is returned.
func FindNetworkStrict(namespace []string) (network *godocker.Network, err error) {
	client, err := Client()
	if err != nil {
		return
	}
	network, err = client.NetworkInfo(Namespace(namespace))
	return
}

// FindNetwork finds a Docker Network by a namespace. If no network if found, NO error is returned.
func FindNetwork(namespace []string) (network *godocker.Network, err error) {
	network, err = FindNetworkStrict(namespace)
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

// AttachNetworkToContainer attaches a network to a container. The network cannot be of driver "overlay".
func AttachNetworkToContainer(networkNamespace []string, containerNamespace []string) (err error) {
	client, err := Client()
	if err != nil {
		return
	}
	network, err := FindNetwork(networkNamespace)
	if err != nil {
		return
	}
	container, err := FindContainerStrict(containerNamespace)
	if err != nil {
		return
	}
	ip, err := FindContainerIP(networkNamespace, containerNamespace)
	if err != nil {
		return
	}
	if ip != "" {
		return
	}
	err = client.ConnectNetwork(network.ID, godocker.NetworkConnectionOptions{
		Container: container.ID,
	})
	if err != nil {
		return
	}
	return
}

// func AttachNetworkToService(network *godocker.Network) (err error) {
// client, err := Client()
// if err != nil {
// return
// }
// options := &ServiceOptions{
// 	Image:     "mesg/daemon",
// 	Namespace: []string{"daemon"},
// 	Ports: []Port{
// 		Port{
// 			Target:    50052,
// 			Published: 50052,
// 		},
// 	},
// 	Mounts: []Mount{
// 		Mount{
// 			Source: "/var/run/docker.sock",
// 			Target: "/var/run/docker.sock",
// 		},
// 		Mount{
// 			Source: viper.GetString(config.MESGPath),
// 			Target: "/mesg",
// 		},
// 	},
// 	NetworksID: []string{network.ID},
// }
// options.merge()

// err = client.UpdateService("ub1mkic25wi90fuiwxkono71u", godocker.UpdateServiceOptions{
// 	Version:     4839,
// 	ServiceSpec: options.CreateServiceOptions.ServiceSpec,
// })
// 	if err != nil {
// 		return
// 	}
// 	return
// }
