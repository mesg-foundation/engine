package daemon

import (
	godocker "github.com/fsouza/go-dockerclient"
	"github.com/mesg-foundation/core/docker"
)

// SharedNetwork returns the shared network
func SharedNetwork() (network *godocker.Network, err error) {
	return docker.FindNetwork(namespaceNetwork())
}
