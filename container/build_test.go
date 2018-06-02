package container

import (
	"testing"

	"github.com/stvp/assert"
)

func TestBuild(t *testing.T) {
	tag, err := Build("test/")
	assert.Nil(t, err)
	assert.NotEqual(t, "", tag)
}
