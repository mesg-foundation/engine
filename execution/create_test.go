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
	require.Nil(t, err)
	require.Equal(t, id, hash.Calculate([]string{
		execution.CreatedAt.UTC().String(),
		execution.Service.Name,
		execution.Task,
		inputs,
	}))
}

func TestCreate(t *testing.T) {
	s := service.Service{
		Name: "TestCreate",
		Tasks: map[string]*service.Task{
			"test": {},
		},
	}
	var inputs map[string]interface{}
	exec, err := Create(&s, "test", inputs)
	require.Nil(t, err)
	require.Equal(t, exec.Service, &s)
	require.Equal(t, exec.Inputs, inputs)
	require.Equal(t, exec.Task, "test")
	require.Equal(t, pendingExecutions[exec.ID], exec)
}

func TestCreateInvalidTask(t *testing.T) {
	s := service.Service{
		Name: "TestCreateInvalidTask",
		Tasks: map[string]*service.Task{
			"test": {},
		},
	}
	var inputs map[string]interface{}
	_, err := Create(&s, "testinvalid", inputs)
	require.NotNil(t, err)
	_, notFound := err.(*service.TaskNotFoundError)
	require.True(t, notFound)
}

func TestCreateInvalidInputs(t *testing.T) {
	s := service.Service{
		Name: "TestCreateInvalidInputs",
		Tasks: map[string]*service.Task{
			"test": {
				Inputs: map[string]*service.Parameter{
					"foo": {
						Type: "String",
					},
				},
			},
		},
	}
	var inputs map[string]interface{}
	_, err := Create(&s, "test", inputs)
	require.NotNil(t, err)
	_, invalid := err.(*service.InvalidTaskInputError)
	require.True(t, invalid)
}
