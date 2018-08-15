package utils

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
	require.Contains(t, cannotReachTheCore, errorMessage(testCoreConnectionErr))
	require.Contains(t, startCore, errorMessage(testCoreConnectionErr))

	require.Contains(t, cannotReachDocker, errorMessage(testDockerConnectionErr))
	require.Contains(t, installDocker, errorMessage(testDockerConnectionErr))

	require.Contains(t, "errorX", errorMessage(errors.New("errorX")))
}
