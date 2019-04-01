package api

import (
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

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
	at.cm.On("Status", mock.Anything).Once().Return(container.RUNNING, nil)

	id, err := a.ExecuteTask("2", "4", map[string]interface{}{}, []string{})
	require.NoError(t, err)
	require.NotNil(t, id)

	at.cm.AssertExpectations(t)
}

func TestExecuteTaskWithInvalidTaskName(t *testing.T) {
	a, at := newTesting(t)
	defer at.close()

	require.NoError(t, at.serviceDB.Save(testService))
	at.cm.On("Status", mock.Anything).Once().Return(container.RUNNING, nil)

	_, err := a.ExecuteTask("2", "2a", map[string]interface{}{}, []string{})
	require.Error(t, err)
}

func TestExecuteTaskForNotRunningService(t *testing.T) {
	a, at := newTesting(t)
	defer at.close()

	require.NoError(t, at.serviceDB.Save(testService))
	at.cm.On("Status", mock.Anything).Once().Return(container.STOPPED, nil)

	_, err := a.ExecuteTask("2", "4", map[string]interface{}{}, []string{})
	_, notRunningError := err.(*NotRunningServiceError)
	require.True(t, notRunningError)
}
