package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	_, err := New()
	require.NoError(t, err)
}

func TestDefaultConfig(t *testing.T) {
	home, _ := homedir.Dir()
	c, err := New()
	require.NoError(t, err)
	require.Equal(t, ":50052", c.Server.Address)
	require.Equal(t, "text", c.Log.Format)
	require.Equal(t, "info", c.Log.Level)
	require.Equal(t, false, c.Log.ForceColors)
	require.Equal(t, filepath.Join(home, ".mesg"), c.Path)
	require.Equal(t, filepath.Join("database", "executions", executionDBVersion), c.Database.ExecutionRelativePath)
	require.Equal(t, "engine", c.Name)
}
func TestEnv(t *testing.T) {
	os.Setenv(envPathKey, "tempPath")
	defer os.Unsetenv(envPathKey)
	os.Setenv(envNameKey, "name")
	defer os.Unsetenv(envNameKey)
	c, err := New()
	require.NoError(t, err)
	require.Equal(t, "tempPath", c.Path)
	require.Equal(t, "name", c.Name)
}

func TestLoadFromFile(t *testing.T) {
	tempPath, _ := ioutil.TempDir("", "TestConfigLoad")
	defer os.RemoveAll(tempPath)
	os.Setenv(envPathKey, tempPath)
	defer os.Unsetenv(envPathKey)
	ioutil.WriteFile(filepath.Join(tempPath, defaultConfigFileName), []byte(`name: shouldNotReadThis
path: shouldNotReadThis
server:
  address: :50050
tendermint:
  config:
    consensus:
      timeoutcommit: 1m6s
`), 0644)

	// load
	c, err := New()
	require.NoError(t, err)
	require.Equal(t, ":50050", c.Server.Address)
	require.Equal(t, 66*time.Second, c.Tendermint.Config.Consensus.TimeoutCommit)
	require.Equal(t, "tcp://0.0.0.0:26657", c.Tendermint.Config.RPC.ListenAddress)
	require.Equal(t, tempPath, c.Path)
	require.Equal(t, "engine", c.Name)
}
