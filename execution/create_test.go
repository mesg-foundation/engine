package execution

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/utils/hash"
	"github.com/stretchr/testify/require"
)

func TestGenerateID(t *testing.T) {
	inputs := "{}"
	var i map[string]interface{}
	json.Unmarshal([]byte(inputs), &i)
	execution := Execution{
		CreatedAt: time.Now(),
		Service:   &service.Service{Name: "TestGenerateID"},
		Task:      "test",
		Inputs:    i,
	}
	id, err := generateID(&execution)
	require.NoError(t, err)
	require.Equal(t, id, hash.Calculate([]string{
		execution.CreatedAt.UTC().String(),
		execution.Service.Name,
		execution.Task,
		inputs,
	}))
}

func TestCreate(t *testing.T) {
	s, _ := service.FromService(&service.Service{
		Name: "TestCreate",
		Tasks: []*service.Task{
			{Key: "test"},
		},
	})
	var inputs map[string]interface{}
	exec, err := Create(s, "test", inputs, []string{})
	require.NoError(t, err)
	require.Equal(t, exec.Service, s)
	require.Equal(t, exec.Inputs, inputs)
	require.Equal(t, exec.Task, "test")
	require.Equal(t, pendingExecutions[exec.ID], exec)
}

func TestCreateInvalidTask(t *testing.T) {
	var (
		taskKey        = "test"
		invalidTaskKey = "testinvalid"
		serviceName    = "TestCreateInvalidTask"
	)
	s, _ := service.FromService(&service.Service{
		Name: serviceName,
		Tasks: []*service.Task{
			{Key: taskKey},
		},
	})
	var inputs map[string]interface{}
	_, err := Create(s, invalidTaskKey, inputs, []string{})
	require.Error(t, err)
	notFoundErr, ok := err.(*service.TaskNotFoundError)
	require.True(t, ok)
	require.Equal(t, invalidTaskKey, notFoundErr.TaskKey)
	require.Equal(t, serviceName, notFoundErr.ServiceName)
}
func TestCreateInvalidInputs(t *testing.T) {
	var (
		taskKey     = "test"
		serviceName = "TestCreateInvalidInputs"
	)
	s, _ := service.FromService(&service.Service{
		Name: serviceName,
		Tasks: []*service.Task{
			{
				Key: taskKey,
				Inputs: []*service.Parameter{
					{
						Key:  "foo",
						Type: "String",
					},
				},
			},
		},
	})
	var inputs map[string]interface{}
	_, err := Create(s, taskKey, inputs, []string{})
	require.Error(t, err)
	invalidErr, ok := err.(*service.InvalidTaskInputError)
	require.True(t, ok)
	require.Equal(t, taskKey, invalidErr.TaskKey)
	require.Equal(t, serviceName, invalidErr.ServiceName)
}
