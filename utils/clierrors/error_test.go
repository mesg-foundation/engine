package clierrors

import (
	"errors"
	"testing"

	"github.com/docker/docker/client"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var testCoreConnectionErr = status.Error(codes.Unavailable, "test")
var testDockerConnectionErr = client.ErrorConnectionFailed("test")

func TestCoreConnectionError(t *testing.T) {
	require.True(t, coreConnectionError(testCoreConnectionErr))
	require.False(t, coreConnectionError(nil))
	require.False(t, coreConnectionError(errors.New("test")))
}

func TestDockerDaemonError(t *testing.T) {
	require.True(t, dockerDaemonError(testDockerConnectionErr))
	require.False(t, dockerDaemonError(nil))
	require.False(t, dockerDaemonError(errors.New("test")))
}

func TestErrorMessage(t *testing.T) {
	require.Contains(t, ErrorMessage(testCoreConnectionErr), cannotReachTheCore)
	require.Contains(t, ErrorMessage(testCoreConnectionErr), startCore)

	require.Contains(t, ErrorMessage(testDockerConnectionErr), cannotReachDocker)
	require.Contains(t, ErrorMessage(testDockerConnectionErr), installDocker)

	require.Contains(t, ErrorMessage(errors.New("errorX")), "errorX")
}
