// Copyright 2018 MESG Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
