package core

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/docker/pkg/archive"
	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/container/dockertest"
	"github.com/mesg-foundation/core/database"
	"github.com/stretchr/testify/require"
)

var (
	eventServicePath = filepath.Join("..", "..", "..", "service-test", "event")
	taskServicePath  = filepath.Join("..", "..", "..", "service-test", "task")
)

const (
	servicedbname = "service.db.test"
	execdbname    = "exec.db.test"
)

func newServerWithContainer(t *testing.T, c container.Container) (*Server, func()) {
	db, err := database.NewServiceDB(servicedbname)
	require.NoError(t, err)

	execDB, err := database.NewExecutionDB(execdbname)
	require.NoError(t, err)

	a, err := api.New(db, execDB, api.ContainerOption(c))
	require.NoError(t, err)

	server := NewServer(a)

	closer := func() {
		db.Close()
		execDB.Close()
		os.RemoveAll(servicedbname)
		os.RemoveAll(execdbname)
	}
	return server, closer
}

func newServer(t *testing.T) (*Server, func()) {
	c, err := container.New()
	require.NoError(t, err)
	return newServerWithContainer(t, c)
}

func newServerAndDockerTest(t *testing.T) (*Server, *dockertest.Testing, func()) {
	dt := dockertest.New()
	c, err := container.New(container.ClientOption(dt.Client()))
	require.NoError(t, err)
	s, closer := newServerWithContainer(t, c)
	return s, dt, closer
}

func serviceTar(t *testing.T, path string) io.Reader {
	reader, err := archive.TarWithOptions(path, &archive.TarOptions{
		Compression: archive.Gzip,
	})
	require.NoError(t, err)
	return reader
}
