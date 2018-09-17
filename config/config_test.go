package config

import (
	"os"
	"path/filepath"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/require"
)

func TestDefaultValue(t *testing.T) {
	home, err := homedir.Dir()
	require.NoError(t, err)

	c, err := New()
	require.NoError(t, err)
	require.Equal(t, ":50052", c.Server.Address)
	require.Equal(t, "localhost:50052", c.Client.Address)
	require.Equal(t, filepath.Join(home, ".mesg", "db"), c.Database.Path)
	require.Equal(t, "text", c.Log.Format)
	require.Equal(t, "info", c.Log.Level)
}

func TestGlobal(t *testing.T) {
	c, err := Global()
	require.NoError(t, err)
	require.NotNil(t, c)
}

func TestLoad(t *testing.T) {
	snapsnot := map[string]string{
		"MESG_SERVER_ADDRESS": "",
		"MESG_CLIENT_ADDRESS": "",
		"MESG_DATABASE_PATH":  "",
		"MESG_LOG_FORMAT":     "",
		"MESG_LOG_LEVEL":      "",
	}
	for key := range snapsnot {
		snapsnot[key] = os.Getenv(key)
	}
	defer func() {
		for key, value := range snapsnot {
			os.Setenv(key, value)
		}
	}()

	os.Setenv("MESG_SERVER_ADDRESS", "test_server_address")
	os.Setenv("MESG_CLIENT_ADDRESS", "test_client_address")
	os.Setenv("MESG_DATABASE_PATH", "test_database_path")
	os.Setenv("MESG_LOG_FORMAT", "test_log_format")
	os.Setenv("MESG_LOG_LEVEL", "test_log_level")

	c, err := New()
	require.NoError(t, err)
	c.LoadFromEnv()

	require.Equal(t, "test_server_address", c.Server.Address)
	require.Equal(t, "test_client_address", c.Client.Address)
	require.Equal(t, "test_database_path", c.Database.Path)
	require.Equal(t, "test_log_format", c.Log.Format)
	require.Equal(t, "test_log_level", c.Log.Level)
}

func TestValidate(t *testing.T) {
	var newConfig = func() *Config {
		c, err := New()
		require.NoError(t, err)
		return c
	}

	c := newConfig()
	require.NoError(t, c.Validate())

	c = newConfig()
	c.Log.Format = "wrongValue"
	require.Error(t, c.Validate())

	c = newConfig()
	c.Log.Level = "wrongValue"
	require.Error(t, c.Validate())

	c = newConfig()
	c.Server.Plugins = append(c.Server.Plugins, plugin{})
	require.Error(t, c.Validate())

	c = newConfig()
	c.Server.Plugins = append(c.Server.Plugins, plugin{})
	require.Error(t, c.Validate())

	c = newConfig()
	c.Server.Plugins = append(c.Server.Plugins, plugin{
		Name: "plugin",
		Path: "/tmp",
		ID:   "1",
	})
	require.Error(t, c.Validate())
}
