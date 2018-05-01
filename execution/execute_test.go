package execution

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/types"
	"github.com/stvp/assert"
)

func TestExecute(t *testing.T) {
	s := service.Service{
		Name: "TestExecute",
		Tasks: map[string]*types.ProtoTask{
			"test": &types.ProtoTask{},
		},
	}
	var inputs interface{}
	execution, _ := Create(&s, "test", inputs)
	reply, err := execution.Execute()
	assert.Nil(t, err)
	assert.NotNil(t, reply)
	assert.Equal(t, reply.ExecutionID, execution.ID)
}

func TestExecuteNotPending(t *testing.T) {
	execution := Execution{
		ID: "TestExecuteNotPending",
	}
	_, err := execution.Execute()
	assert.NotNil(t, err)
}
