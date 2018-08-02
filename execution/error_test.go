package execution

import (
	"testing"

	"github.com/stvp/assert"
)

func TestNotInQueueError(t *testing.T) {
	e := NotInQueueError{"test", "queueName"}
	assert.Contains(t, "Execution 'test' not found in queue 'queueName'", e.Error())
}
