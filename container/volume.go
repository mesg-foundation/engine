package container

import (
	"context"
)

// DeleteVolume deletes a Docker Volume by name.
func (c *DockerContainer) DeleteVolume(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.callTimeout)
	defer cancel()
	return c.client.VolumeRemove(ctx, name, false)
}
