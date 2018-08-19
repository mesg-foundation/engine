package execution

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestComplete(t *testing.T) {
	s := service.Service{
		Name: "TestComplete",
		Tasks: map[string]*service.Task{
			"test": {
				Outputs: map[string]*service.Output{
					"output": {},
				},
			},
		},
	}
	var inputs map[string]interface{}
	execution, _ := Create(&s, "test", inputs, []string{})
	execution.Execute()
	var outputs map[string]interface{}
	err := execution.Complete("output", outputs)
	require.Nil(t, err)
	require.Equal(t, execution.Output, "output")
	require.Equal(t, execution.OutputData, outputs)
	require.True(t, execution.ExecutionDuration > 0)
}

func TestCompleteNotFound(t *testing.T) {
	s := service.Service{
		Name: "TestCompleteNotFound",
		Tasks: map[string]*service.Task{
			"test": {},
		},
	}
	var inputs map[string]interface{}
	execution, _ := Create(&s, "test", inputs, []string{})
	execution.Execute()
	var outputs map[string]interface{}
	err := execution.Complete("output", outputs)
	require.NotNil(t, err)
	x, missingOutputError := err.(*service.OutputNotFoundError)
	require.True(t, missingOutputError)
	require.Equal(t, "output", x.OutputKey)
}

func TestCompleteInvalidOutputs(t *testing.T) {
	s := service.Service{
		Name: "TestCompleteInvalidOutputs",
		Tasks: map[string]*service.Task{
			"test": {
				Outputs: map[string]*service.Output{
					"output": {
						Data: map[string]*service.Parameter{
							"foo": {},
						},
					},
				},
			},
		},
	}
	var inputs map[string]interface{}
	execution, _ := Create(&s, "test", inputs, []string{})
	execution.Execute()
	var outputs map[string]interface{}
	err := execution.Complete("output", outputs)
	require.NotNil(t, err)
	x, invalidOutputError := err.(*service.InvalidOutputDataError)
	require.True(t, invalidOutputError)
	require.Equal(t, "output", x.Key)
}

func TestCompleteNotProcessed(t *testing.T) {
	s := service.Service{
		Name: "TestCompleteNotProcessed",
		Tasks: map[string]*service.Task{
			"test": {
				Outputs: map[string]*service.Output{
					"output": {},
				},
			},
		},
	}
	var inputs map[string]interface{}
	execution, _ := Create(&s, "test", inputs, []string{})
	var outputs map[string]interface{}
	err := execution.Complete("output", outputs)
	require.NotNil(t, err)
	x, notInQueueError := err.(*NotInQueueError)
	require.True(t, notInQueueError)
	require.Equal(t, "inProgress", x.Queue)
}
