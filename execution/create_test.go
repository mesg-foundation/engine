package execution

import (
	"testing"
	"time"

	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/types"
	"github.com/mesg-foundation/core/utils/hash"
	"github.com/stvp/assert"
)

func TestGenerateID(t *testing.T) {
	execution := Execution{
		CreatedAt: time.Now(),
		Service:   &service.Service{Name: "TestGenerateID"},
		Task:      "test",
	}
	id := generateID(&execution)
	assert.Equal(t, id, hash.Calculate([]string{
		execution.CreatedAt.UTC().String(),
		execution.Service.Name,
		execution.Task,
	}))
}

func TestTaskExists(t *testing.T) {
	s := service.Service{
		Tasks: map[string]*types.ProtoTask{
			"test": &types.ProtoTask{},
		},
	}
	exisits := taskExists(&s, "test")
	assert.True(t, exisits)
}

func TestTaskNotExists(t *testing.T) {
	s := service.Service{
		Tasks: map[string]*types.ProtoTask{
			"test": &types.ProtoTask{},
		},
	}
	exisits := taskExists(&s, "testnotexists")
	assert.False(t, exisits)
}

func TestTaskExistsOnOtherKey(t *testing.T) {
	s := service.Service{
		Tasks: map[string]*types.ProtoTask{
			"test":  &types.ProtoTask{},
			"test2": &types.ProtoTask{},
		},
	}
	exisits := taskExists(&s, "test2")
	assert.True(t, exisits)
}

func TestCreate(t *testing.T) {
	s := service.Service{
		Name: "TestCreate",
		Tasks: map[string]*types.ProtoTask{
			"test": &types.ProtoTask{},
		},
	}
	var inputs interface{}
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
		Tasks: map[string]*types.ProtoTask{
			"test": &types.ProtoTask{},
		},
	}
	var inputs interface{}
	_, err := Create(&s, "testinvalid", inputs)
	assert.NotNil(t, err)
}
