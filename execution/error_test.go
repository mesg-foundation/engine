package execution

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNotInQueueError(t *testing.T) {
	e := NotInQueueError{"test", "queueName"}
	require.Contains(t, e.Error(), "Execution 'test' not found in queue 'queueName'")
}
