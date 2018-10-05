package execution

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestComplete(t *testing.T) {
	s, _ := service.FromService(&service.Service{
		Name: "TestComplete",
		Tasks: []*service.Task{
			{
				Key: "test",
				Outputs: []*service.Output{
					{
						Key: "output",
					},
				},
			},
		},
	})
	var inputs map[string]interface{}
	execution, _ := Create(s, "test", inputs, []string{})
	execution.Execute()
	var outputs map[string]interface{}
	err := execution.Complete("output", outputs)
	require.Nil(t, err)
	require.Equal(t, execution.Output, "output")
	require.Equal(t, execution.OutputData, outputs)
	require.True(t, execution.ExecutionDuration > 0)
}

func TestCompleteNotProcessed(t *testing.T) {
	s, _ := service.FromService(&service.Service{
		Name: "TestCompleteNotProcessed",
		Tasks: []*service.Task{
			{
				Key: "test",
				Outputs: []*service.Output{
					{
						Key: "output",
					},
				},
			},
		},
	})
	var inputs map[string]interface{}
	execution, _ := Create(s, "test", inputs, []string{})
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
	s, _ := service.FromService(&service.Service{
		Name: serviceName,
		Tasks: []*service.Task{
			{
				Key: taskKey,
			},
		},
	})
	var inputs map[string]interface{}
	execution, _ := Create(s, taskKey, inputs, []string{})
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
		taskKey     = "test"
		outputKey   = "output"
		serviceName = "TestCompleteInvalidOutputs"
	)
	s, _ := service.FromService(&service.Service{
		Name: serviceName,
		Tasks: []*service.Task{
			{
				Key: taskKey,
				Outputs: []*service.Output{
					{
						Key: outputKey,
						Data: []*service.Parameter{
							{Key: "foo"},
						},
					},
				},
			},
		},
	})
	var inputs map[string]interface{}
	execution, _ := Create(s, taskKey, inputs, []string{})
	execution.Execute()
	var outputs map[string]interface{}
	err := execution.Complete(outputKey, outputs)
	require.NotNil(t, err)
	invalidErr, ok := err.(*service.InvalidTaskOutputError)
	require.True(t, ok)
	require.Equal(t, taskKey, invalidErr.TaskKey)
	require.Equal(t, outputKey, invalidErr.TaskOutputKey)
	require.Equal(t, serviceName, invalidErr.ServiceName)
}
