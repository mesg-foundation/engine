// +build integration

package container

import (
	"testing"

	"github.com/stvp/assert"
)

func TestIntegrationBuild(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	tag, err := c.Build("test/")
	assert.Nil(t, err)
	assert.NotEqual(t, "", tag)
}

func TestIntegrationBuildNotWorking(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	tag, err := c.Build("test-not-valid/")
	assert.NotNil(t, err)
	assert.Equal(t, "", tag)
}

func TestIntegrationBuildWrongPath(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	_, err = c.Build("testss/")
	assert.NotNil(t, err)
}
