package execution

import (
	"errors"
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	var (
		s      = &service.Service{}
		inputs = make(map[string]interface{})
		tags   = make([]string, 0)

		exec = New(s, "id", "foo", inputs, tags)
	)

	assert.Equal(t, s, exec.Service)
	assert.Equal(t, "id", exec.EventID)
	assert.Equal(t, "foo", exec.TaskKey)
	assert.Equal(t, inputs, exec.Inputs)
	assert.Equal(t, tags, exec.Tags)
	assert.Equal(t, Created, exec.Status)
	assert.False(t, exec.CreatedAt.IsZero())
}

func TestExecute(t *testing.T) {
	exec := New(nil, "id", "foo", nil, nil)
	assert.NoError(t, exec.Execute())
	assert.Equal(t, InProgress, exec.Status)
	assert.False(t, exec.ExecutedAt.IsZero())
	assert.Error(t, exec.Execute())
}

func TestComplte(t *testing.T) {
	var (
		s = &service.Service{
			Tasks: []*service.Task{{
				Key: "foo",
				Outputs: []*service.Output{
					{
						Key: "output",
					},
				},
			}},
		}

		data = map[string]interface{}{"some": "data"}
		exec = New(s, "id", "foo", nil, nil)
	)

	exec.Execute()
	assert.NoError(t, exec.Complete("output", data))
	assert.Equal(t, Completed, exec.Status)
	assert.Equal(t, data, exec.OutputData)
	assert.Equal(t, "output", exec.OutputKey)
	assert.Error(t, exec.Complete("output", data))
}

func TestFailed(t *testing.T) {
	exec := New(nil, "id", "foo", nil, nil)
	exec.Execute()
	assert.NoError(t, exec.Failed(errors.New("error")))
	assert.Equal(t, Failed, exec.Status)
	assert.Equal(t, "error", exec.Error)
	assert.True(t, exec.ExecutionDuration > 0)
	assert.Error(t, exec.Failed(errors.New("error")))
}
