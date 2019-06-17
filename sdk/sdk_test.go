package sdk

import (
	"os"
	"testing"

	"github.com/mesg-foundation/core/container/mocks"
	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/service/manager/dockermanager"
	"github.com/stretchr/testify/require"
)

const (
	servicedbname  = "service.db.test"
	instancedbname = "instance.db.test"
	execdbname     = "exec.db.test"
)

type apiTesting struct {
	*testing.T
	serviceDB     *database.LevelDBServiceDB
	instanceDB    *database.LevelDBInstanceDB
	executionDB   *database.LevelDBExecutionDB
	containerMock *mocks.Container
}

func (t *apiTesting) close() {
	require.NoError(t, t.serviceDB.Close())
	require.NoError(t, t.executionDB.Close())
	require.NoError(t, t.instanceDB.Close())
	require.NoError(t, os.RemoveAll(servicedbname))
	require.NoError(t, os.RemoveAll(execdbname))
	require.NoError(t, os.RemoveAll(instancedbname))
}

func newTesting(t *testing.T) (*SDK, *apiTesting) {
	containerMock := &mocks.Container{}
	m := dockermanager.New(containerMock) // TODO(ilgooz): create mocks from manager.Manager and use instead.

	db, err := database.NewServiceDB(servicedbname)
	require.NoError(t, err)

	instanceDB, err := database.NewInstanceDB(instancedbname)
	require.NoError(t, err)

	execDB, err := database.NewExecutionDB(execdbname)
	require.NoError(t, err)

	a := New(m, containerMock, db, instanceDB, execDB)

	return a, &apiTesting{
		T:             t,
		serviceDB:     db,
		instanceDB:    instanceDB,
		executionDB:   execDB,
		containerMock: containerMock,
	}
}
