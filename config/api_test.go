package config

import (
	"path/filepath"
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
		APIServerAddress:       ":50052",
		APIServerSocket:        "/mesg/server.sock",
		APIClientTarget:        viper.GetString(APIServerAddress),
		APIServiceSocketPath:   filepath.Join(viper.GetString(MESGPath), "server.sock"),
		APIServiceTargetPath:   "/mesg/server.sock",
		APIServiceTargetSocket: "unix://" + viper.GetString(APIServiceTargetPath),
		LogFormat:              "text",
		LogLevel:               "info",
		ServicePathHost:        filepath.Join(viper.GetString(MESGPath), "services"),
		ServicePathDocker:      filepath.Join("/mesg", "services"),
	}
	for key, defaultValue := range defaults {
		assertViperDefault(t, key, defaultValue)
	}

	// Override by ENV when testing, so only test the image name
	require.Contains(t, viper.GetString(CoreImage), "mesg/core:")
}
