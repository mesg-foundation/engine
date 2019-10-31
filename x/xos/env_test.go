package xos

import (
	"testing"
)

func TestEnvMergeSlices(t *testing.T) {
	values := [][]string{
		{"a=1", "b=2"},
		{"c=3", "a=2"},
	}
	env := EnvMergeSlices(values...)
	for i, v := range []string{"a=2", "b=2", "c=3"} {
		if env[i] != v {
			t.Errorf("env slice dosen't contain %s", v)
		}
	}
}
