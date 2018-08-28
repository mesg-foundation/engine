package api

import (
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/container/dockertest"
	"github.com/stretchr/testify/require"
)

func newAPIAndDockerTest(t *testing.T) (*API, *dockertest.Testing) {
	dt := dockertest.New()

	container, err := container.New(container.ClientOption(dt.Client()))
	require.Nil(t, err)

	a, err := New(ContainerOption(container))
	require.Nil(t, err)

	return a, dt
}
