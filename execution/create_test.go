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
	exec, err := Create(&s, "test", inputs, []string{})
	require.Nil(t, err)
	require.Equal(t, exec.Service, &s)
	require.Equal(t, exec.Inputs, inputs)
	require.Equal(t, exec.Task, "test")
	require.Equal(t, pendingExecutions[exec.ID], exec)
}
