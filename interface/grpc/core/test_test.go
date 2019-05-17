package core

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/mesg-foundation/core/service/manager/dockermanager"

	"github.com/docker/docker/pkg/archive"
	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/container"
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
	m := dockermanager.New(c) // TODO(ilgooz): create mocks from manager.Manager and use instead.

	db, err := database.NewServiceDB(servicedbname)
	require.NoError(t, err)

	execDB, err := database.NewExecutionDB(execdbname)
	require.NoError(t, err)

	a := api.New(m, c, db, execDB)

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

func serviceTar(t *testing.T, path string) io.Reader {
	reader, err := archive.TarWithOptions(path, &archive.TarOptions{
		Compression: archive.Gzip,
	})
	require.NoError(t, err)
	return reader
}
