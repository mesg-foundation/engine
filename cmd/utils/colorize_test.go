package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestColorizeJSON(t *testing.T) {
	jsonStringified := `{"key1":"value1","key2":"value2"}`
	colorized := ColorizeJSON(jsonStringified)
	require.Equal(t, "\x1b[36mkey1\x1b[0m = value1, \x1b[36mkey2\x1b[0m = value2", colorized)
}
