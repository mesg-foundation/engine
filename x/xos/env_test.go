package xos

import (
	"testing"

	"github.com/mesg-foundation/engine/x/xstrings"
)

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
