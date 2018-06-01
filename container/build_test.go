package container

import (
	"testing"

	"github.com/stvp/assert"
)

func TestBuild(t *testing.T) {
	tag, err := Build("test/", []string{"testbuild"})
	assert.Nil(t, err)
	assert.NotEqual(t, "", tag)
}
