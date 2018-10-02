// +build integration

package service

import (
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/stretchr/testify/require"
)

func newIntegrationContainer(t *testing.T) *container.DockerContainer {
	c, err := container.New()
	require.NoError(t, err)
	return c
}
