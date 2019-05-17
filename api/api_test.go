package api

import (
	"os"
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/container/mocks"
	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/service/manager/dockermanager"
	"github.com/mesg-foundation/core/utils/hash"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	servicedbname = "service.db.test"
	execdbname    = "exec.db.test"
)

type apiTesting struct {
	*testing.T
	serviceDB     *database.LevelDBServiceDB
	executionDB   *database.LevelDBExecutionDB
	containerMock *mocks.Container
}

func (t *apiTesting) close() {
	require.NoError(t, t.serviceDB.Close())
	require.NoError(t, t.executionDB.Close())
	require.NoError(t, os.RemoveAll(servicedbname))
	require.NoError(t, os.RemoveAll(execdbname))
}

func newTesting(t *testing.T) (*API, *apiTesting) {
	containerMock := &mocks.Container{}
	m := dockermanager.New(containerMock) // TODO(ilgooz): create mocks from manager.Manager and use instead.

	db, err := database.NewServiceDB(servicedbname)
	require.NoError(t, err)

	execDB, err := database.NewExecutionDB(execdbname)
	require.NoError(t, err)

	a := New(m, containerMock, db, execDB)

	return a, &apiTesting{
		T:             t,
		serviceDB:     db,
		executionDB:   execDB,
		containerMock: containerMock,
	}
}

var testService = &service.Service{
	Name: "1",
	Sid:  "2",
	Hash: "33",
	Tasks: []*service.Task{
		{Key: "4"},
	},
	Dependencies: []*service.Dependency{
		{Key: "5"},
	},
}

func TestNotRunningServiceError(t *testing.T) {
	e := NotRunningServiceError{ServiceID: "test"}
	require.Equal(t, `Service "test" is not running`, e.Error())
}

func TestExecuteTask(t *testing.T) {
	a, at := newTesting(t)
	defer at.close()

	// TODO(ilgooz): use api.Deploy() instead of manually saving the service
	// and do the same improvement in the similar places.
	// in order to do this, create a testing helper to build service tarballs
	// from yml definitions on the fly .
	require.NoError(t, at.serviceDB.Save(testService))
	at.containerMock.On("Status", mock.Anything).Once().Return(container.RUNNING, nil)

	id, err := a.ExecuteTask("2", "4", map[string]interface{}{}, []string{})
	require.NoError(t, err)
	require.NotNil(t, id)

	at.containerMock.AssertExpectations(t)
}

func TestExecuteTaskWithInvalidTaskName(t *testing.T) {
	a, at := newTesting(t)
	defer at.close()

	require.NoError(t, at.serviceDB.Save(testService))
	at.containerMock.On("Status", mock.Anything).Once().Return(container.RUNNING, nil)

	_, err := a.ExecuteTask("2", "2a", map[string]interface{}{}, []string{})
	require.Error(t, err)
}

func TestExecuteTaskForNotRunningService(t *testing.T) {
	a, at := newTesting(t)
	defer at.close()

	require.NoError(t, at.serviceDB.Save(testService))
	at.containerMock.On("Status", mock.Anything).Once().Return(container.STOPPED, nil)

	_, err := a.ExecuteTask("2", "4", map[string]interface{}{}, []string{})
	_, notRunningError := err.(*NotRunningServiceError)
	require.True(t, notRunningError)
}

func TestEventSubTopic(t *testing.T) {
	serviceHash := "1"
	require.Equal(t, eventSubTopic(serviceHash), hash.Calculate([]string{serviceHash, eventTopic}))
}

func TestExecutionSubTopic(t *testing.T) {
	serviceHash := "1"
	require.Equal(t, executionSubTopic(serviceHash), hash.Calculate([]string{serviceHash, executionTopic}))
}
