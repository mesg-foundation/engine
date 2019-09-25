package core

import (
	"os"
	"testing"

	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/database"
	"github.com/mesg-foundation/engine/sdk"
	"github.com/stretchr/testify/require"
)

const (
	instancedbname = "instance.db.test"
	execdbname     = "exec.db.test"
	processdbname  = "process.db.test"
)

func newServerWithContainer(t *testing.T, c container.Container) (*Server, func()) {
	instanceDB, err := database.NewInstanceDB(instancedbname)
	require.NoError(t, err)

	execDB, err := database.NewExecutionDB(execdbname)
	require.NoError(t, err)

	processDB, err := database.NewProcessDB(processdbname)
	require.NoError(t, err)

	a := sdk.New(nil, nil, nil, c, instanceDB, execDB, processDB, "", "")

	server := NewServer(a)

	closer := func() {
		instanceDB.Close()
		execDB.Close()
		processDB.Close()
		os.RemoveAll(instancedbname)
		os.RemoveAll(execdbname)
		os.RemoveAll(processdbname)
	}
	return server, closer
}

func newServer(t *testing.T) (*Server, func()) {
	c, err := container.New("enginetest")
	require.NoError(t, err)
	return newServerWithContainer(t, c)
}
