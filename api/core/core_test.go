package core

import (
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/container/dockertest"
	"github.com/mesg-foundation/core/mesg"
	"github.com/stretchr/testify/require"
)

func newServer(t *testing.T) *Server {
	container, err := container.New()
	require.Nil(t, err)

	m, err := mesg.New(mesg.DockerClientOption(container))
	require.Nil(t, err)

	server, err := NewServer(MESGOption(m))
	require.Nil(t, err)

	return server
}

func newServerAndDockerTest(t *testing.T) (*Server, *dockertest.Testing) {
	dt := dockertest.New()

	container, err := container.New(container.ClientOption(dt.Client()))
	require.Nil(t, err)

	m, err := mesg.New(mesg.DockerClientOption(container))
	require.Nil(t, err)

	server, err := NewServer(MESGOption(m))
	require.Nil(t, err)

	return server, dt
}
