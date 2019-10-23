package config

import (
	"os"
	"path/filepath"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	_, err := New()
	require.NoError(t, err)
}

func TestDefaultValue(t *testing.T) {
	home, _ := homedir.Dir()
	c, err := Default()
	require.NoError(t, err)
	require.Equal(t, ":50052", c.Server.Address)
	require.Equal(t, "text", c.Log.Format)
	require.Equal(t, "info", c.Log.Level)
	require.Equal(t, false, c.Log.ForceColors)
	require.Equal(t, filepath.Join(home, ".mesg"), c.Path)
	require.Equal(t, filepath.Join("database", "executions", executionDBVersion), c.Database.ExecutionRelativePath)
	require.Equal(t, "engine", c.Name)
}

func TestLoad(t *testing.T) {
	snapsnot := map[string]string{
		"MESG_SERVER_ADDRESS":                 "",
		"MESG_LOG_FORMAT":                     "",
		"MESG_LOG_LEVEL":                      "",
		"MESG_LOG_FORCECOLORS":                "",
		"MESG_TENDERMINT_P2P_PERSISTENTPEERS": "",
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
	os.Setenv("MESG_LOG_FORMAT", "test_log_format")
	os.Setenv("MESG_LOG_LEVEL", "test_log_level")
	os.Setenv("MESG_LOG_FORCECOLORS", "true")
	os.Setenv("MESG_TENDERMINT_P2P_PERSISTENTPEERS", "localhost")

	c, _ := Default()
	c.Load()
	require.Equal(t, "test_server_address", c.Server.Address)
	require.Equal(t, "test_log_format", c.Log.Format)
	require.Equal(t, "test_log_level", c.Log.Level)
	require.Equal(t, true, c.Log.ForceColors)
	require.Equal(t, "localhost", c.Tendermint.P2P.PersistentPeers)
}
