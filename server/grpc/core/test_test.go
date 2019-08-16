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
	servicedbname  = "service.db.test"
	instancedbname = "instance.db.test"
	execdbname     = "exec.db.test"
	workflowdbname = "workflow.db.test"
)

func newServerWithContainer(t *testing.T, c container.Container) (*Server, func()) {
	db, err := database.NewServiceDB(servicedbname)
	require.NoError(t, err)

	instanceDB, err := database.NewInstanceDB(instancedbname)
	require.NoError(t, err)

	execDB, err := database.NewExecutionDB(execdbname)
	require.NoError(t, err)

	workflowDB, err := database.NewWorkflowDB(workflowdbname)
	require.NoError(t, err)

	a := sdk.New(c, db, instanceDB, execDB, workflowDB, "", "")

	server := NewServer(a)

	closer := func() {
		db.Close()
		execDB.Close()
		os.RemoveAll(servicedbname)
		os.RemoveAll(instancedbname)
		os.RemoveAll(execdbname)
	}
	return server, closer
}

func newServer(t *testing.T) (*Server, func()) {
	c, err := container.New("enginetest")
	require.NoError(t, err)
	return newServerWithContainer(t, c)
}
