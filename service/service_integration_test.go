// +build integration

package service

import (
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/stretchr/testify/require"
)

func newIntegrationContainer(t *testing.T, options ...container.Option) container.Container {
	c, err := container.New(options...)
	require.NoError(t, err)
	return c
}
