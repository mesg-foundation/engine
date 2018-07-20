package execution

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/utils/hash"
	"github.com/stvp/assert"
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
	assert.Nil(t, err)
	assert.Equal(t, id, hash.Calculate([]string{
		execution.CreatedAt.UTC().String(),
		execution.Service.Name,
		execution.Task,
		inputs,
	}))
}

// func TestTaskExists(t *testing.T) {
// 	s := service.Service{
// 		Tasks: map[string]*service.Task{
// 			"test": &service.Task{},
// 		},
// 	}
// 	exisits := taskExists(&s, "test")
// 	assert.True(t, exisits)
// }

// func TestTaskNotExists(t *testing.T) {
// 	s := service.Service{
// 		Tasks: map[string]*service.Task{
// 			"test": &service.Task{},
// 		},
// 	}
// 	exisits := taskExists(&s, "testnotexists")
// 	assert.False(t, exisits)
// }

// func TestTaskExistsOnOtherKey(t *testing.T) {
// 	s := service.Service{
// 		Tasks: map[string]*service.Task{
// 			"test":  &service.Task{},
// 			"test2": &service.Task{},
// 		},
// 	}
// 	exisits := taskExists(&s, "test2")
// 	assert.True(t, exisits)
// }

func TestCreate(t *testing.T) {
	s := service.Service{
		Name: "TestCreate",
		Tasks: map[string]*service.Task{
			"test": &service.Task{},
		},
	}
	var inputs map[string]interface{}
	exec, err := Create(&s, "test", inputs)
	assert.Nil(t, err)
	assert.Equal(t, exec.Service, &s)
	assert.Equal(t, exec.Inputs, inputs)
	assert.Equal(t, exec.Task, "test")
	assert.Equal(t, pendingExecutions[exec.ID], exec)
}

func TestCreateInvalidTask(t *testing.T) {
	s := service.Service{
		Name: "TestCreateInvalidTask",
		Tasks: map[string]*service.Task{
			"test": &service.Task{},
		},
	}
	var inputs map[string]interface{}
	_, err := Create(&s, "testinvalid", inputs)
	assert.NotNil(t, err)
}

func TestCreateInvalidInputs(t *testing.T) {
	s := service.Service{
		Name: "TestCreateInvalidInputs",
		Tasks: map[string]*service.Task{
			"test": &service.Task{
				Inputs: map[string]*service.Parameter{
					"foo": &service.Parameter{
						Type: "String",
					},
				},
			},
		},
	}
	var inputs map[string]interface{}
	_, err := Create(&s, "test", inputs)
	assert.Contains(t, "Invalid inputs", err.Error())
	assert.NotNil(t, err)
}
