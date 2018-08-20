// +build integration

package container

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationBuild(t *testing.T) {
	c, err := New()
	require.Nil(t, err)
	tag, err := c.Build("test/")
	require.Nil(t, err)
	require.NotEqual(t, "", tag)
}

func TestIntegrationBuildNotWorking(t *testing.T) {
	c, err := New()
	require.Nil(t, err)
	tag, err := c.Build("test-not-valid/")
	require.NotNil(t, err)
	require.Equal(t, "", tag)
}

func TestIntegrationBuildWrongPath(t *testing.T) {
	c, err := New()
	require.Nil(t, err)
	_, err = c.Build("testss/")
	require.NotNil(t, err)
}
