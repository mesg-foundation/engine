package service

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/docker/pkg/archive"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/sdk"
	"github.com/mesg-foundation/core/service/manager/dockermanager"
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

func newServer(t *testing.T) (*Server, func()) {
	db, err := database.NewServiceDB(servicedbname)
	require.NoError(t, err)

	execDB, err := database.NewExecutionDB(execdbname)
	require.NoError(t, err)

	c, err := container.New()
	require.NoError(t, err)

	m := dockermanager.New(c) // TODO(ilgooz): create mocks from manager.Manager and use instead.

	a := sdk.New(m, c, db, execDB)
	server := NewServer(a)

	closer := func() {
		require.NoError(t, db.Close())
		require.NoError(t, execDB.Close())
		require.NoError(t, os.RemoveAll(servicedbname))
		require.NoError(t, os.RemoveAll(execdbname))
	}

	return server, closer
}

func serviceTar(t *testing.T, path string) io.Reader {
	reader, err := archive.TarWithOptions(path, &archive.TarOptions{
		Compression: archive.Gzip,
	})
	require.NoError(t, err)
	return reader
}

func jsonMarshal(t *testing.T, data interface{}) string {
	bytes, err := json.Marshal(data)
	require.NoError(t, err)
	return string(bytes)
}
