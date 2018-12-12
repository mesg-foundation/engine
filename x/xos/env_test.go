package xos

import (
	"os"
	"testing"

	"github.com/mesg-foundation/core/x/xstrings"
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

func TestMapToEnv(t *testing.T) {
	env := MapToEnv(map[string]string{
		"a": "1",
		"b": "2",
	})
	for _, v := range []string{"a=1", "b=2"} {
		if !xstrings.SliceContains(env, v) {
			t.Errorf("envs dosen't contain %s", v)
		}
	}
}
func TestMergeMapEnvs(t *testing.T) {
	envs := []map[string]string{
		{
			"a": "1",
			"b": "2",
		},
		{
			"a": "2",
			"c": "3",
		},
	}
	env := MergeMapEnvs(envs...)
	for k, v := range map[string]string{"a": "2", "b": "2", "c": "3"} {
		if env[k] != v {
			t.Errorf("envs dosen't contain %s=%s", k, v)
		}
	}
}
