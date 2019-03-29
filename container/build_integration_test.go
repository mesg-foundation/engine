// +build integration

package container

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationBuild(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	tag, err := c.Build("test/", "test", "x1")
	require.NoError(t, err)
	require.NotEqual(t, "", tag)
	require.Equal(t, "test:x1", tag)
}

func TestIntegrationBuildNotWorking(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	_, err = c.Build("test-not-valid/", "test", "x2")
	require.Error(t, err)
}

func TestIntegrationBuildWrongPath(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	_, err = c.Build("testss/", "test", "x3")
	require.Error(t, err)
}
