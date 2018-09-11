package xcolor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNextColor(t *testing.T) {
	for _, color := range colors {
		require.Equal(t, color, NextColor())
	}
	require.Equal(t, colors[0], NextColor())
}
