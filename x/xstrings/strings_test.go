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

func TestAppendSpace(t *testing.T) {
	for _, tt := range []struct {
		s        string
		n        int
		expected string
	}{
		{"a", -1, "a"},
		{"a", 0, "a"},
		{"a", 1, "a "},
		{"a", 2, "a  "},
	} {
		if got := AppendSpace(tt.s, tt.n); got != tt.expected {
			t.Errorf("%v append space %d - got %s, want %s", tt.s, tt.n, got, tt.expected)
		}
	}
}
