package service

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/docker/pkg/archive"
	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/database"
	"github.com/stretchr/testify/require"
)

var (
	eventServicePath = filepath.Join("..", "..", "..", "service-test", "event")
	taskServicePath  = filepath.Join("..", "..", "..", "service-test", "task")
)

func newServer(t *testing.T) (*Server, func()) {
	db, err := database.NewServiceDB("db.test")
	require.NoError(t, err)

	a, err := api.New(db)
	require.NoError(t, err)

	server, err := NewServer(APIOption(a))
	require.NoError(t, err)

	closer := func() {
		db.Close()
		os.RemoveAll("db.test")
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
