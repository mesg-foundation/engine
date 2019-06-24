package service

import (
	"os"
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/sdk"
	"github.com/stretchr/testify/require"
)

const (
	servicedbname  = "service.db.test"
	instancedbname = "instance.db.test"
	execdbname     = "exec.db.test"
)

func newServer(t *testing.T) (*Server, func()) {
	db, err := database.NewServiceDB(servicedbname)
	require.NoError(t, err)

	instanceDB, err := database.NewInstanceDB(instancedbname)
	require.NoError(t, err)

	execDB, err := database.NewExecutionDB(execdbname)
	require.NoError(t, err)

	c, err := container.New()
	require.NoError(t, err)

	a := sdk.New(c, db, instanceDB, execDB)
	server := NewServer(a)

	closer := func() {
		require.NoError(t, db.Close())
		require.NoError(t, execDB.Close())
		require.NoError(t, os.RemoveAll(servicedbname))
		require.NoError(t, os.RemoveAll(instancedbname))
		require.NoError(t, os.RemoveAll(execdbname))
	}

	return server, closer
}
