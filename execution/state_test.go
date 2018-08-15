package execution

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestMoveFromPendingToInProgress(t *testing.T) {
	s := service.Service{
		Name: "TestMoveFromPendingToInProgress",
		Tasks: map[string]*service.Task{
			"test": {},
		},
	}
	var inputs map[string]interface{}
	exec, _ := Create(&s, "test", inputs)
	err := exec.moveFromPendingToInProgress()
	require.Equal(t, inProgressExecutions[exec.ID], exec)
	require.Nil(t, pendingExecutions[exec.ID])
	require.Nil(t, err)
}

func TestMoveFromPendingToInProgressNonExistingTask(t *testing.T) {
	exec := Execution{ID: "test"}
	err := exec.moveFromPendingToInProgress()
	require.NotNil(t, err)
	require.Nil(t, pendingExecutions[exec.ID])
}

func TestMoveFromInProgressToCompleted(t *testing.T) {
	s := service.Service{
		Name: "TestMoveFromInProgressToCompleted",
		Tasks: map[string]*service.Task{
			"test": {},
		},
	}
	var inputs map[string]interface{}
	exec, _ := Create(&s, "test", inputs)
	exec.moveFromPendingToInProgress()
	err := exec.moveFromInProgressToProcessed()
	require.Equal(t, processedExecutions[exec.ID], exec)
	require.Nil(t, inProgressExecutions[exec.ID])
	require.Nil(t, err)
}

func TestMoveFromInProgressToCompletedNonExistingTask(t *testing.T) {
	s := service.Service{
		Name: "TestMoveFromInProgressToCompletedNonExistingTask",
		Tasks: map[string]*service.Task{
			"test": {},
		},
	}
	var inputs map[string]interface{}
	exec, _ := Create(&s, "test", inputs)
	err := exec.moveFromInProgressToProcessed()
	require.NotNil(t, err)
	require.Nil(t, inProgressExecutions[exec.ID])
}

func TestInProgress(t *testing.T) {
	inProgressExecutions["foo"] = &Execution{ID: "TestInProgress"}
	require.NotNil(t, InProgress("foo"))
	require.Equal(t, "TestInProgress", InProgress("foo").ID)
	require.Nil(t, InProgress("bar"))
}
