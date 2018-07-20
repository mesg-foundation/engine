package execution

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

func TestExecute(t *testing.T) {
	s := service.Service{
		Name: "TestExecute",
		Tasks: map[string]*service.Task{
			"test": &service.Task{},
		},
	}
	var inputs map[string]interface{}
	execution, _ := Create(&s, "test", inputs)
	err := execution.Execute()
	assert.Nil(t, err)
}

func TestExecuteNotPending(t *testing.T) {
	execution := Execution{
		ID: "TestExecuteNotPending",
	}
	err := execution.Execute()
	assert.NotNil(t, err)
}
