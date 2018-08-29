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

func TestCompleteNotFound(t *testing.T) {
	var (
		taskKey     = "test"
		outputKey   = "output"
		serviceName = "TestCompleteNotFound"
	)
	s := service.Service{
		Name: serviceName,
		Tasks: map[string]*service.Task{
			taskKey: {},
		},
	}
	var inputs map[string]interface{}
	execution, _ := Create(&s, taskKey, inputs, []string{})
	execution.Execute()
	var outputs map[string]interface{}
	err := execution.Complete(outputKey, outputs)
	require.NotNil(t, err)
	notFoundErr, ok := err.(*service.TaskOutputNotFoundError)
	require.True(t, ok)
	require.Equal(t, taskKey, notFoundErr.TaskKey)
	require.Equal(t, outputKey, notFoundErr.TaskOutputKey)
	require.Equal(t, serviceName, notFoundErr.ServiceName)
}

func TestCompleteInvalidOutputs(t *testing.T) {
	var (
		taskKey   = "test"
		outputKey = "output"
	)
	s := service.Service{
		Name: "TestCompleteInvalidOutputs",
		Tasks: map[string]*service.Task{
			taskKey: {
				Outputs: map[string]*service.Output{
					outputKey: {
						Data: map[string]*service.Parameter{
							"foo": {},
						},
					},
				},
			},
		},
	}
	var inputs map[string]interface{}
	execution, _ := Create(&s, taskKey, inputs, []string{})
	execution.Execute()
	var outputs map[string]interface{}
	err := execution.Complete(outputKey, outputs)
	require.NotNil(t, err)
	invalidErr, ok := err.(*service.InvalidTaskOutputError)
	require.True(t, ok)
	require.Equal(t, taskKey, invalidErr.TaskKey)
	require.Equal(t, outputKey, invalidErr.TaskOutputKey)
}
