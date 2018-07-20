package execution

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

func TestMissingOutputError(t *testing.T) {
	e := MissingOutputError{Service: &service.Service{Name: "test"}, Output: "test"}
	assert.Equal(t, "Output test doesn't exists in service test", e.Error())
}

func TestInvalidOutputError(t *testing.T) {
	e := InvalidOutputError{Service: &service.Service{Name: "test"}, Warnings: []*service.ParameterWarning{}}
	assert.Contains(t, "Invalid result: ", e.Error())
}

func TestNotInQueueError(t *testing.T) {
	e := NotInQueueError{"test", "queueName"}
	assert.Contains(t, "Execution test not in queueName queue", e.Error())
}
