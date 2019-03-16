package container

import (
	"context"
)

// DeleteVolume deletes a Docker Volume by name.
func (c *DockerContainer) DeleteVolume(name string) error {
	return c.client.VolumeRemove(context.Background(), name, false)
}
