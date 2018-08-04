package config

import (
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stvp/assert"
)

func assertViperDefault(t *testing.T, key string, expected string) {
	value := viper.GetString(key)
	assert.Equal(t, expected, value, "Wrong default for key "+key)
}

func TestAPIDefault(t *testing.T) {
	defaults := map[string]string{
		APIServerAddress:       ":50052",
		APIServerSocket:        "/mesg/server.sock",
		APIClientTarget:        viper.GetString(APIServerAddress),
		APIServiceSocketPath:   filepath.Join(viper.GetString(MESGPath), "server.sock"),
		APIServiceTargetPath:   "/mesg/server.sock",
		APIServiceTargetSocket: "unix://" + viper.GetString(APIServiceTargetPath),
		ServicePathHost:        filepath.Join(viper.GetString(MESGPath), "services"),
		ServicePathDocker:      filepath.Join("/mesg", "services"),
	}
	for key, defaultValue := range defaults {
		assertViperDefault(t, key, defaultValue)
	}

	// Override by ENV when testing, so only test the image name
	assert.Contains(t, "mesg/core:", viper.GetString(CoreImage))
}
