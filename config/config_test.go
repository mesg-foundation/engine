package config

import (
	"testing"

	"github.com/stvp/assert"
)

func TestToEnv(t *testing.T) {
	assert.Equal(t, envPrefix+envSeparator+"FOO"+envSeparator+"BAR", ToEnv("foo.bar"))
}
