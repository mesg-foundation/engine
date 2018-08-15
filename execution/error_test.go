package execution

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNotInQueueError(t *testing.T) {
	e := NotInQueueError{"test", "queueName"}
	require.Contains(t, "Execution 'test' not found in queue 'queueName'", e.Error())
}
