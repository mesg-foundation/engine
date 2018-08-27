package xos

import (
	"os"
	"testing"
)

func TestGetenvDefault(t *testing.T) {
	os.Setenv("__TEST_KEY__", "1")
	for _, tt := range []struct {
		key      string
		fallback string
		expected string
	}{
		{"__TEST_KEY__", "0", "1"},
		{"__TEST_KE1Y__", "0", "0"},
	} {
		if got := GetenvDefault(tt.key, tt.fallback); got != tt.expected {
			t.Errorf("GetenvDefault(%q, %q) got %s, want %s", tt.key, tt.fallback, got, tt.expected)
		}
	}
}
