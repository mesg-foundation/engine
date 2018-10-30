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

const testdbname = "db.test"

func newServer(t *testing.T) (*Server, func()) {
	container, err := container.New()
	require.NoError(t, err)

	db, err := database.NewServiceDB(testdbname)
	require.NoError(t, err)

	execDB, err := database.NewExecutionDB("execution" + testdbname)
	require.NoError(t, err)

	a, err := api.New(db, execDB, api.ContainerOption(container))
	require.NoError(t, err)

	server := NewServer(a)

	closer := func() {
		db.Close()
		execDB.Close()
		os.RemoveAll(testdbname)
	}
	return server, closer
}

func newServerAndDockerTest(t *testing.T) (*Server, *dockertest.Testing, func()) {
	dt := dockertest.New()

	container, err := container.New(container.ClientOption(dt.Client()))
	require.NoError(t, err)

	db, err := database.NewServiceDB("qeffeq" + testdbname)
	require.NoError(t, err)

	execDB, err := database.NewExecutionDB("dqjbdq" + testdbname)
	require.NoError(t, err)

	a, err := api.New(db, execDB, api.ContainerOption(container))
	require.NoError(t, err)

	server := NewServer(a)

	closer := func() {
		db.Close()
		execDB.Close()
		os.RemoveAll(testdbname)
	}
	return server, dt, closer
}

func serviceTar(t *testing.T, path string) io.Reader {
	reader, err := archive.TarWithOptions(path, &archive.TarOptions{
		Compression: archive.Gzip,
	})
	require.NoError(t, err)
	return reader
}
