package utils

import (
	"errors"
	"testing"

	"github.com/docker/docker/client"
	"github.com/stvp/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var testCoreConnectionErr = status.Error(codes.Unavailable, "test")
var testDockerConnectionErr = client.ErrorConnectionFailed("test")

func TestCoreConnectionError(t *testing.T) {
	assert.True(t, coreConnectionError(testCoreConnectionErr))
	assert.False(t, coreConnectionError(nil))
	assert.False(t, coreConnectionError(errors.New("test")))
}

func TestDockerDaemonError(t *testing.T) {
	assert.True(t, dockerDaemonError(testDockerConnectionErr))
	assert.False(t, dockerDaemonError(nil))
	assert.False(t, dockerDaemonError(errors.New("test")))
}

func TestErrorMessage(t *testing.T) {
	assert.Contains(t, cannotReachTheCore, errorMessage(testCoreConnectionErr))
	assert.Contains(t, startCore, errorMessage(testCoreConnectionErr))

	assert.Contains(t, cannotReachDocker, errorMessage(testDockerConnectionErr))
	assert.Contains(t, installDocker, errorMessage(testDockerConnectionErr))

	assert.Contains(t, "errorX", errorMessage(errors.New("errorX")))
}
