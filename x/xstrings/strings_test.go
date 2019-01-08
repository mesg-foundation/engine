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

package xstrings

import "testing"

func TestSliceContains(t *testing.T) {
	for _, tt := range []struct {
		s        []string
		e        string
		expected bool
	}{
		{[]string{"a"}, "a", true},
		{[]string{"a"}, "b", false},
	} {
		if got := SliceContains(tt.s, tt.e); got != tt.expected {
			t.Errorf("%v slice contains %s - got %t, want %t", tt.s, tt.e, got, tt.expected)
		}
	}
}

func TestFindLongest(t *testing.T) {
	for _, tt := range []struct {
		s        []string
		expected int
	}{
		{[]string{"a"}, 1},
		{[]string{"a", "aa"}, 2},
	} {
		if got := FindLongest(tt.s); got != tt.expected {
			t.Errorf("%v slice find longetst - got %d, want %d", tt.s, got, tt.expected)
		}
	}
}

func TestSliceIndex(t *testing.T) {
	for _, tt := range []struct {
		s        []string
		e        string
		expected int
	}{
		{[]string{"a"}, "b", -1},
		{[]string{"a", "b"}, "a", 0},
		{[]string{"a", "b"}, "b", 1},
	} {
		if got := SliceIndex(tt.s, tt.e); got != tt.expected {
			t.Errorf("%v slice index %s - got %d, want %d", tt.s, tt.e, got, tt.expected)
		}
	}
}
