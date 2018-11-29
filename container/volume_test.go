package container

import (
	"testing"

	"github.com/docker/docker/api/types"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateVolume(t *testing.T) {
	var (
		c, m   = newTesting(t)
		name   = "1"
		volume = types.Volume{Name: "2"}
	)

	m.On("VolumeCreate", mock.Anything, volumetypes.VolumeCreateBody{Name: name}).Once().Return(volume, nil)

	volume1, err := c.CreateVolume(name)
	require.NoError(t, err)
	require.Equal(t, volume, volume1)

	m.AssertExpectations(t)
}

func TestDeleteVolume(t *testing.T) {
	var (
		c, m = newTesting(t)
		name = "1"
	)

	m.On("VolumeRemove", mock.Anything, name, false).Once().Return(nil)
	require.NoError(t, c.DeleteVolume(name))

	m.AssertExpectations(t)
}
