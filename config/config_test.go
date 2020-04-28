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
	require.Equal(t, "engine", c.Name)
	require.Equal(t, "engine", c.Account.Name)
	require.Equal(t, "pass", c.Account.Password)
	require.Equal(t, uint32(0), c.Account.Number)
	require.Equal(t, uint32(0), c.Account.Index)
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
	tempPath, _ := ioutil.TempDir("", "TestLoadFromFile")
	defer os.RemoveAll(tempPath)
	os.Setenv(envPathKey, tempPath)
	defer os.Unsetenv(envPathKey)

	t.Run("key does not exist", func(t *testing.T) {
		ioutil.WriteFile(filepath.Join(tempPath, defaultConfigFileName), []byte(`foo: bar`), 0644)
		_, err := New()
		require.Error(t, err)
	})
	t.Run("load", func(t *testing.T) {
		ioutil.WriteFile(filepath.Join(tempPath, defaultConfigFileName), []byte(`server:
  address: :50050
log:
  forcecolors: true
account:
  mnemonic: glimpse upon body vast economy give taxi yellow rabbit come click ranch chronic hammer sport near rotate charge lumber chicken cloud base thing forum
cosmos:
  mingasprices: 2.0019294mesg
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
		require.Equal(t, "2.0019294mesg", c.Cosmos.MinGasPrices)
	})
}
