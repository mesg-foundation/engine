package xos

import (
	"os"
	"testing"

	"github.com/mesg-foundation/engine/x/xstrings"
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

func TestEnvMapToSlice(t *testing.T) {
	env := EnvMapToSlice(map[string]string{
		"a": "1",
		"b": "2",
	})
	for _, v := range []string{"a=1", "b=2"} {
		if !xstrings.SliceContains(env, v) {
			t.Errorf("env slice dosen't contain %s", v)
		}
	}
}

func TestEnvMapToString(t *testing.T) {
	got := EnvMapToString(map[string]string{
		"b": "2",
		"a": "1",
	})
	want := "a=1;b=2"
	if got != want {
		t.Errorf("invalid env string - got %s, want %s", got, want)
	}
}

func TestEnvSliceToMap(t *testing.T) {
	env := EnvSliceToMap([]string{"a=1", "b=2"})
	for k, v := range map[string]string{"a": "1", "b": "2"} {
		if env[k] != v {
			t.Errorf("env map dosen't contain %s=%v", k, v)
		}
	}
}

func TestEnvMergeMaps(t *testing.T) {
	values := []map[string]string{
		{
			"a": "1",
			"b": "2",
		},
		{
			"a": "2",
			"c": "3",
		},
	}
	env := EnvMergeMaps(values...)
	for k, v := range map[string]string{"a": "2", "b": "2", "c": "3"} {
		if env[k] != v {
			t.Errorf("env map dosen't contain %s=%s", k, v)
		}
	}
}
func TestEnvMergeSlices(t *testing.T) {
	values := [][]string{
		{"a=1", "b=2"},
		{"a=2", "c=3"},
	}
	env := EnvMergeSlices(values...)
	for i, v := range []string{"a=1", "b=2", "a=2", "c=3"} {
		if env[i] != v {
			t.Errorf("env slice dosen't contain %s", v)
		}
	}
}
