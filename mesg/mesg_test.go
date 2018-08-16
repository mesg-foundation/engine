package mesg

import (
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/container/dockertest"
	"github.com/stretchr/testify/require"
)

func newMESGAndDockerTest(t *testing.T) (*MESG, *dockertest.Testing) {
	dt := dockertest.New()

	container, err := container.New(container.ClientOption(dt.Client()))
	require.Nil(t, err)

	m, err := New(DockerClientOption(container))
	require.Nil(t, err)

	return m, dt
}
