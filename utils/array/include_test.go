package array

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIncludedIn(t *testing.T) {
	require.False(t, IncludedIn([]string{}, ""))
	require.True(t, IncludedIn([]string{""}, ""))
	require.False(t, IncludedIn([]string{"a"}, ""))
	require.True(t, IncludedIn([]string{"a"}, "a"))
	require.False(t, IncludedIn([]string{""}, "a"))
	require.True(t, IncludedIn([]string{"a", "b"}, "a"))
	require.True(t, IncludedIn([]string{"a", "b"}, "b"))
	require.False(t, IncludedIn([]string{"a", "b"}, "c"))
}
