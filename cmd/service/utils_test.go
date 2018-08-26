package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultPath(t *testing.T) {
	require.Equal(t, defaultPath([]string{}), "./")
	require.Equal(t, defaultPath([]string{"foo"}), "foo")
	require.Equal(t, defaultPath([]string{"foo", "bar"}), "foo")
}
