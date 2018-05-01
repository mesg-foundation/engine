package execution

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/types"
	"github.com/stvp/assert"
)

func TestComplete(t *testing.T) {
	s := service.Service{
		Name: "TestComplete",
		Tasks: map[string]*types.ProtoTask{
			"test": &types.ProtoTask{},
		},
	}
	var inputs interface{}
	execution, _ := Create(&s, "test", inputs)
	execution.Execute()
	var outputs interface{}
	err := execution.Complete("output", outputs)
	assert.Nil(t, err)
	assert.Equal(t, execution.Output, "output")
	assert.Equal(t, execution.OutputData, outputs)
	assert.True(t, execution.ExecutionDuration > 0)
}

func TestCompleteNotProcessed(t *testing.T) {
	s := service.Service{
		Name: "TestCompleteNotProcessed",
		Tasks: map[string]*types.ProtoTask{
			"test": &types.ProtoTask{},
		},
	}
	var inputs interface{}
	execution, _ := Create(&s, "test", inputs)
	var outputs interface{}
	err := execution.Complete("output", outputs)
	assert.NotNil(t, err)
}
