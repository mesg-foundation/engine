package dockertest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestNotFoundErr makes sure NotFoundErr to implement docker client's notFound interface.
func TestNotFoundErr(t *testing.T) {
	err := NotFoundErr{}
	require.True(t, err.NotFound())
	require.Equal(t, "not found", err.Error())
}
