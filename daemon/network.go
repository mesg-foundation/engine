package daemon

import (
	godocker "github.com/fsouza/go-dockerclient"
	"github.com/mesg-foundation/core/docker"
)

// IP returns the IP of the daemon in the shared network
func IP() (daemonIP string, err error) {
	return docker.FindServiceIP(NamespaceNetwork(), Namespace())
}

// SharedNetwork returns the shared network
func SharedNetwork() (network *godocker.Network, err error) {
	return docker.FindNetwork(NamespaceNetwork())
}
