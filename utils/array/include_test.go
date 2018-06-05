package array

import (
	"testing"

	"github.com/stvp/assert"
)

func TestIncludedIn(t *testing.T) {
	assert.False(t, IncludedIn([]string{}, ""))
	assert.True(t, IncludedIn([]string{""}, ""))
	assert.False(t, IncludedIn([]string{"a"}, ""))
	assert.True(t, IncludedIn([]string{"a"}, "a"))
	assert.False(t, IncludedIn([]string{""}, "a"))
	assert.True(t, IncludedIn([]string{"a", "b"}, "a"))
	assert.True(t, IncludedIn([]string{"a", "b"}, "b"))
	assert.False(t, IncludedIn([]string{"a", "b"}, "c"))
}
