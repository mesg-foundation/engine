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
