package daemon

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogs(t *testing.T) {
	startForTest()
	reader, err := Logs()
	require.NoError(t, err)
	require.NotNil(t, reader)
}
