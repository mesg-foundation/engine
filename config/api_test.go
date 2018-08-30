package config

import (
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func assertViperDefault(t *testing.T, defaults map[string]string) {
	for key, defaultValue := range defaults {
		value := viper.GetString(key)
		require.Equal(t, defaultValue, value, "Wrong default for key "+key)
	}
}

func TestAPIDefault(t *testing.T) {
	assertViperDefault(t, map[string]string{
		APIServerAddress:  ":50052",
		ServicePathHost:   filepath.Join(viper.GetString(MESGPath), "services"),
		ServicePathDocker: filepath.Join("/mesg", "services"),
	})

	// Override by ENV when testing, so only test the image name
	require.Contains(t, viper.GetString(CoreImage), "mesg/core:")
}
