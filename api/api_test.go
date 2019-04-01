package api

import (
	"os"
	"testing"

	"github.com/mesg-foundation/core/container/mocks"
	"github.com/mesg-foundation/core/database"
	"github.com/stretchr/testify/require"
)

const (
	servicedbname = "service.db.test"
	execdbname    = "exec.db.test"
)

type apiTesting struct {
	*testing.T
	serviceDB   *database.LevelDBServiceDB
	executionDB *database.LevelDBExecutionDB
	cm          *mocks.Container
}

func (t *apiTesting) close() {
	require.NoError(t, t.serviceDB.Close())
	require.NoError(t, t.executionDB.Close())
	require.NoError(t, os.RemoveAll(servicedbname))
	require.NoError(t, os.RemoveAll(execdbname))
}

func newTesting(t *testing.T) (*API, *apiTesting) {
	cm := &mocks.Container{}

	db, err := database.NewServiceDB(servicedbname)
	require.NoError(t, err)

	execDB, err := database.NewExecutionDB(execdbname)
	require.NoError(t, err)

	a, err := New(db, execDB, ContainerOption(cm))
	require.NoError(t, err)

	return a, &apiTesting{
		T:           t,
		serviceDB:   db,
		executionDB: execDB,
		cm:          cm,
	}
}
