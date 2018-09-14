package execution

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestExecute(t *testing.T) {
	s, _ := service.FromService(&service.Service{
		Name: "TestExecute",
		Tasks: []*service.Task{
			{Key: "test"},
		},
	})
	var inputs map[string]interface{}
	execution, _ := Create(s, "test", inputs, []string{})
	err := execution.Execute()
	require.Nil(t, err)
}

func TestExecuteNotPending(t *testing.T) {
	execution := Execution{
		ID: "TestExecuteNotPending",
	}
	err := execution.Execute()
	require.NotNil(t, err)
}
