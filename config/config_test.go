package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultValue(t *testing.T) {
	tests := []struct {
		entry        func() *Entry
		defaultValue string
		contains     bool
	}{
		{
			entry:        APIPort,
			defaultValue: "50052",
		},
		{
			entry:        APIAddress,
			defaultValue: "",
		},
		{
			entry:        LogFormat,
			defaultValue: "text",
		},
		{
			entry:        LogLevel,
			defaultValue: "info",
		},
		{
			entry:        CoreImage,
			defaultValue: "mesg/core:",
			contains:     true,
		},
	}
	for _, test := range tests {
		value, err := test.entry().GetValue()
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
		entry        func() *Entry
		value        string
		shouldFail   bool
		errorMessage string
	}{
		{
			entry:      LogFormat,
			value:      "text",
			shouldFail: false,
		},
		{
			entry:      LogFormat,
			value:      "json",
			shouldFail: false,
		},
		{
			entry:        LogFormat,
			value:        "wrongFormat",
			shouldFail:   true,
			errorMessage: `Value "wrongFormat" is not an allowed`,
		},
		{
			entry:      LogLevel,
			value:      "warning",
			shouldFail: false,
		},
		{
			entry:        LogLevel,
			value:        "wrongLevel",
			shouldFail:   true,
			errorMessage: `Value "wrongLevel" is not a valid log level`,
		},
	}
	for _, test := range tests {
		err := test.entry().SetValue(test.value)
		if test.shouldFail {
			require.EqualError(t, err, test.errorMessage)
		} else {
			require.NoError(t, err)
		}
	}
}
