package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func assertViperDefault(t *testing.T, key string, expected string) {
	value := viper.GetString(key)
	require.Equal(t, expected, value, "Wrong default for key "+key)
}

func TestAPIDefault(t *testing.T) {
	defaults := map[string]string{
		Path:            "/mesg",
		APIPort:    "50052",
		APIAddress: "",
		LogFormat:  "text",
		LogLevel:   "info",
	}
	for key, defaultValue := range defaults {
		assertViperDefault(t, key, defaultValue)
	}

	// Override by ENV when testing, so only test the image name
	require.Contains(t, viper.GetString(CoreImage), "mesg/core:")
}
