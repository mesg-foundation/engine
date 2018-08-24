package core

import (
	"io"
	"testing"

	"github.com/docker/docker/pkg/archive"
	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/container/dockertest"
	"github.com/stretchr/testify/require"
)

func newServer(t *testing.T) *Server {
	container, err := container.New()
	require.NoError(t, err)

	a, err := api.New(api.ContainerOption(container))
	require.NoError(t, err)

	server, err := NewServer(APIOption(a))
	require.NoError(t, err)

	return server
}

func newServerAndDockerTest(t *testing.T) (*Server, *dockertest.Testing) {
	dt := dockertest.New()

	container, err := container.New(container.ClientOption(dt.Client()))
	require.NoError(t, err)

	a, err := api.New(api.ContainerOption(container))
	require.NoError(t, err)

	server, err := NewServer(APIOption(a))
	require.NoError(t, err)

	return server, dt
}

func serviceTar(t *testing.T, path string) io.Reader {
	reader, err := archive.TarWithOptions(path, &archive.TarOptions{
		Compression: archive.Gzip,
	})
	require.NoError(t, err)
	return reader
}
