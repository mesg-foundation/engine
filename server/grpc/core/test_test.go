package core

import (
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/database"
	"github.com/mesg-foundation/engine/database/store"
	"github.com/mesg-foundation/engine/sdk"
	"github.com/stretchr/testify/require"
)

const (
	servicedbname  = "service.db.test"
	instancedbname = "instance.db.test"
	execdbname     = "exec.db.test"
	processdbname  = "process.db.test"
)

func newServerWithContainer(t *testing.T, c container.Container) (*Server, func()) {
	serviceStore, err := store.NewLevelDBStore(servicedbname)
	require.NoError(t, err)
	db := database.NewServiceDB(serviceStore, codec.New())

	instanceDB, err := database.NewInstanceDB(instancedbname)
	require.NoError(t, err)

	execDB, err := database.NewExecutionDB(execdbname)
	require.NoError(t, err)

	processDB, err := database.NewProcessDB(processdbname)
	require.NoError(t, err)

	a := sdk.NewDeprecated(c, db, instanceDB, execDB, processDB, "", "")

	server := NewServer(a)

	closer := func() {
		db.Close()
		instanceDB.Close()
		execDB.Close()
		processDB.Close()
		os.RemoveAll(servicedbname)
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
