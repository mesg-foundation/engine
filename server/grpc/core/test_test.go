package core

import (
	"os"
	"testing"

	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/database"
	"github.com/mesg-foundation/engine/database/store"
	"github.com/mesg-foundation/engine/sdk"
	servicesdk "github.com/mesg-foundation/engine/sdk/service"
	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	servicedbname  = "service.db.test"
	instancedbname = "instance.db.test"
	execdbname     = "exec.db.test"
	workflowdbname = "workflow.db.test"
)

func newServerWithContainer(t *testing.T, c container.Container) (*Server, func()) {
	s, err := leveldb.OpenFile(servicedbname, nil)
	require.NoError(t, err)
	serviceDB := database.NewServiceDB(store.NewLevelDBStore(s))
	serviceSDK := servicesdk.NewDeprecated(c, serviceDB)

	instanceDB, err := database.NewInstanceDB(instancedbname)
	require.NoError(t, err)

	execDB, err := database.NewExecutionDB(execdbname)
	require.NoError(t, err)

	workflowDB, err := database.NewWorkflowDB(workflowdbname)
	require.NoError(t, err)

	a := sdk.New(c, serviceSDK, instanceDB, execDB, workflowDB, "", "")

	server := NewServer(a)

	closer := func() {
		serviceDB.Close()
		instanceDB.Close()
		execDB.Close()
		workflowDB.Close()
		os.RemoveAll(servicedbname)
		os.RemoveAll(instancedbname)
		os.RemoveAll(execdbname)
		os.RemoveAll(workflowdbname)
	}
	return server, closer
}

func newServer(t *testing.T) (*Server, func()) {
	c, err := container.New("enginetest")
	require.NoError(t, err)
	return newServerWithContainer(t, c)
}
