package dockermanager

import (
	"crypto/sha1"
	"encoding/hex"

	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/hash"
)

// DockerManager is responsible for managing MESG Service's Docker Containers
// in Docker Service environment.
type DockerManager struct {
	c container.Container
}

// New returns a new Docker Manager.
func New(c container.Container) *DockerManager {
	return &DockerManager{
		c: c,
	}
}

// serviceNamespace returns the namespace of the service.
func serviceNamespace(hash hash.Hash) []string {
	sum := sha1.Sum(hash)
	return []string{hex.EncodeToString(sum[:])}
}

// dependencyNamespace builds the namespace of a dependency.
func dependencyNamespace(serviceNamespace []string, dependencyKey string) []string {
	return append(serviceNamespace, dependencyKey)
}
