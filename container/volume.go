package container

import (
	"context"

	"github.com/docker/docker/api/types"
	volumetypes "github.com/docker/docker/api/types/volume"
)

// CreateVolume creates a Docker Volume with name.
func (c *DockerContainer) CreateVolume(name string) (types.Volume, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.callTimeout)
	defer cancel()
	return c.client.VolumeCreate(ctx, volumetypes.VolumeCreateBody{
		Name: name,
	})
}

// DeleteVolume deletes a Docker Volume by name.
func (c *DockerContainer) DeleteVolume(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.callTimeout)
	defer cancel()
	return c.client.VolumeRemove(ctx, name, false)
}
