package xreflect

import (
	"testing"
	"time"
)

func TestIsNil(t *testing.T) {
	var tests = []struct {
		o        interface{}
		name     string
		expected bool
	}{
		{nil, "nil", true},
		{(*time.Time)(nil), "nil struct", true},
		{(**time.Time)(nil), "nil nil struct", true},
		{time.Time{}, "struct", false},
		{&time.Time{}, "ptr struct", false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := IsNil(tt.o)
			if actual != tt.expected {
				t.Errorf("(%s): expected %t, actual %t", tt.name, tt.expected, actual)
			}
		})
	}
}
