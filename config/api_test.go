package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultValue(t *testing.T) {
	tests := []struct {
		config       func() *Config
		defaultValue string
		contains     bool
	}{
		{
			config:       APIPort,
			defaultValue: "50052",
		},
		{
			config:       APIAddress,
			defaultValue: "",
		},
		{
			config:       LogFormat,
			defaultValue: "text",
		},
		{
			config:       LogLevel,
			defaultValue: "info",
		},
		{
			config:       CoreImage,
			defaultValue: "mesg/core:",
			contains:     true,
		},
	}
	for _, test := range tests {
		value, err := test.config().GetValue()
		require.NoError(t, err)
		if test.contains {
			require.Contains(t, value, test.defaultValue)
		} else {
			require.Equal(t, test.defaultValue, value)
		}
	}
}

func TestValidation(t *testing.T) {
	tests := []struct {
		config       func() *Config
		value        string
		shouldFail   bool
		errorMessage string
	}{
		{
			config:     LogFormat,
			value:      "text",
			shouldFail: false,
		},
		{
			config:     LogFormat,
			value:      "json",
			shouldFail: false,
		},
		{
			config:       LogFormat,
			value:        "wrongFormat",
			shouldFail:   true,
			errorMessage: `Value "wrongFormat" is not an allowed`,
		},
		{
			config:     LogLevel,
			value:      "warning",
			shouldFail: false,
		},
		{
			config:       LogLevel,
			value:        "wrongLevel",
			shouldFail:   true,
			errorMessage: `Value "wrongLevel" is not a valid log level`,
		},
	}
	for _, test := range tests {
		err := test.config().SetValue(test.value)
		if test.shouldFail {
			require.EqualError(t, err, test.errorMessage)
		} else {
			require.NoError(t, err)
		}
	}
}
