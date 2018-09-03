package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultValue(t *testing.T) {
	tests := []struct {
		setting      func() Setting
		defaultValue string
	}{
		{APIPort, "50052"},
		{APIAddress, ""},
		{LogFormat, "text"},
		{LogLevel, "info"},
	}
	for _, test := range tests {
		require.Equal(t, test.defaultValue, test.setting().GetValue())
	}

	require.Contains(t, CoreImage().GetValue(), "mesg/core:")
}
