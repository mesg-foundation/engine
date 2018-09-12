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
