package container

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDeleteVolume(t *testing.T) {
	var (
		c, m = newTesting(t)
		name = "1"
	)

	m.On("VolumeRemove", mock.Anything, name, false).Once().Return(nil)
	require.NoError(t, c.DeleteVolume(name))

	m.AssertExpectations(t)
}
