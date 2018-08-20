package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToEnv(t *testing.T) {
	require.Equal(t, envPrefix+envSeparator+"FOO"+envSeparator+"BAR", ToEnv("foo.bar"))
}
