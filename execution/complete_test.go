package execution

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

func TestComplete(t *testing.T) {
	s := service.Service{
		Name: "TestComplete",
		Tasks: map[string]*service.Task{
			"test": &service.Task{
				Outputs: map[string]*service.Output{
					"output": &service.Output{},
				},
			},
		},
	}
	var inputs map[string]interface{}
	execution, _ := Create(&s, "test", inputs)
	execution.Execute()
	var outputs map[string]interface{}
	err := execution.Complete("output", outputs)
	assert.Nil(t, err)
	assert.Equal(t, execution.Output, "output")
	assert.Equal(t, execution.OutputData, outputs)
	assert.True(t, execution.ExecutionDuration > 0)
}

func TestCompleteNotFound(t *testing.T) {
	s := service.Service{
		Name: "TestCompleteNotFound",
		Tasks: map[string]*service.Task{
			"test": &service.Task{},
		},
	}
	var inputs map[string]interface{}
	execution, _ := Create(&s, "test", inputs)
	execution.Execute()
	var outputs map[string]interface{}
	err := execution.Complete("output", outputs)
	assert.NotNil(t, err)
	x, missingOutputError := err.(*MissingOutputError)
	assert.True(t, missingOutputError)
	assert.Equal(t, "output", x.Output)
}

func TestCompleteInvalidOutputs(t *testing.T) {
	s := service.Service{
		Name: "TestCompleteInvalidOutputs",
		Tasks: map[string]*service.Task{
			"test": &service.Task{
				Outputs: map[string]*service.Output{
					"output": &service.Output{
						Data: map[string]*service.Parameter{
							"foo": &service.Parameter{},
						},
					},
				},
			},
		},
	}
	var inputs map[string]interface{}
	execution, _ := Create(&s, "test", inputs)
	execution.Execute()
	var outputs map[string]interface{}
	err := execution.Complete("output", outputs)
	assert.NotNil(t, err)
	x, invalidOutputError := err.(*InvalidOutputError)
	assert.True(t, invalidOutputError)
	assert.True(t, len(x.Warnings) > 0)
}

func TestCompleteNotProcessed(t *testing.T) {
	s := service.Service{
		Name: "TestCompleteNotProcessed",
		Tasks: map[string]*service.Task{
			"test": &service.Task{
				Outputs: map[string]*service.Output{
					"output": &service.Output{},
				},
			},
		},
	}
	var inputs map[string]interface{}
	execution, _ := Create(&s, "test", inputs)
	var outputs map[string]interface{}
	err := execution.Complete("output", outputs)
	assert.NotNil(t, err)
	x, notInQueueError := err.(*NotInQueueError)
	assert.True(t, notInQueueError)
	assert.Equal(t, "inProgress", x.Queue)
}
