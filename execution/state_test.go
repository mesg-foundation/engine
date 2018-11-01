package execution

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestMoveFromPendingToInProgress(t *testing.T) {
	s, _ := service.FromService(&service.Service{
		Name: "TestMoveFromPendingToInProgress",
		Tasks: []*service.Task{
			{Key: "test"},
		},
	})
	var inputs map[string]interface{}
	exec, _ := Create(s, "test", inputs, []string{})
	err := exec.moveFromPendingToInProgress()
	require.Equal(t, inProgressExecutions[exec.ID], exec)
	require.Nil(t, pendingExecutions[exec.ID])
	require.NoError(t, err)
}

func TestMoveFromPendingToInProgressNonExistingTask(t *testing.T) {
	exec := Execution{ID: "test"}
	err := exec.moveFromPendingToInProgress()
	require.Error(t, err)
	require.Nil(t, pendingExecutions[exec.ID])
}

func TestDeleteFromInProgressQueue(t *testing.T) {
	s, _ := service.FromService(&service.Service{
		Name: "TestDeleteFromInProgressQueue",
		Tasks: []*service.Task{
			{Key: "test"},
		},
	})
	var inputs map[string]interface{}
	exec, _ := Create(s, "test", inputs, []string{})
	exec.moveFromPendingToInProgress()
	err := exec.deleteFromInProgressQueue()
	require.Nil(t, inProgressExecutions[exec.ID])
	require.NoError(t, err)
}

func TestDeleteFromInProgressQueueNonExistingTask(t *testing.T) {
	s, _ := service.FromService(&service.Service{
		Name: "TestDeleteFromInProgressQueueNonExistingTask",
		Tasks: []*service.Task{
			{Key: "test"},
		},
	})
	var inputs map[string]interface{}
	exec, _ := Create(s, "test", inputs, []string{})
	err := exec.deleteFromInProgressQueue()
	require.Error(t, err)
	require.Nil(t, inProgressExecutions[exec.ID])
}

func TestInProgress(t *testing.T) {
	inProgressExecutions["foo"] = &Execution{ID: "TestInProgress"}
	require.NotNil(t, InProgress("foo"))
	require.Equal(t, "TestInProgress", InProgress("foo").ID)
	require.Nil(t, InProgress("bar"))
}
