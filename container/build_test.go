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

func TestBuildWrongPath(t *testing.T) {
	_, err := Build("testss/")
	assert.NotNil(t, err)
}

func TestBuildHash(t *testing.T) {
	hash, _ := Build("test/")
	assert.Equal(t, hash, "sha256:acf8f4e61c5bea3743c50f3ac96526d75518bd3290b7a8291d32ece917fd82a5")
}
