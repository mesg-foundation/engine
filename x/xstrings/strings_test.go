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
