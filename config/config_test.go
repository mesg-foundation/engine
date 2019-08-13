package config

import (
	"os"
	"path/filepath"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func TestDefaultValue(t *testing.T) {
	home, _ := homedir.Dir()
	c, err := New()
	require.NoError(t, err)
	require.Equal(t, ":50052", c.Server.Address)
	require.Equal(t, "text", c.Log.Format)
	require.Equal(t, "info", c.Log.Level)
	require.Equal(t, false, c.Log.ForceColors)
	require.Equal(t, filepath.Join(home, ".mesg"), c.Path)
	require.Equal(t, filepath.Join("database", "services", serviceDBVersion), c.Database.ServiceRelativePath)
	require.Equal(t, filepath.Join("database", "executions", executionDBVersion), c.Database.ExecutionRelativePath)
	require.Equal(t, "engine", c.Name)
}

func TestGlobal(t *testing.T) {
	c, err := Global()
	require.NoError(t, err)
	require.NotNil(t, c)
}

func TestLoad(t *testing.T) {
	snapsnot := map[string]string{
		"MESG_SERVER_ADDRESS":                 "",
		"MESG_LOG_FORMAT":                     "",
		"MESG_LOG_LEVEL":                      "",
		"MESG_LOG_FORCECOLORS":                "",
		"MESG_TENDERMINT_P2P_PERSISTENTPEERS": "",
		"MESG_COSMOS_VALIDATOR_PUB_KEY":       "",
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
	os.Setenv("MESG_COSMOS_VALIDATOR_PUB_KEY", "0000000000000000000000000000000000000000000000000000000000000001")

	c, _ := New()
	c.Load()
	require.Equal(t, "test_server_address", c.Server.Address)
	require.Equal(t, "test_log_format", c.Log.Format)
	require.Equal(t, "test_log_level", c.Log.Level)
	require.Equal(t, true, c.Log.ForceColors)
	require.Equal(t, "localhost", c.Tendermint.P2P.PersistentPeers)
	require.Equal(t, "PubKeyEd25519{0000000000000000000000000000000000000000000000000000000000000001}", ed25519.PubKeyEd25519(c.Cosmos.ValidatorPubKey).String())
}

func TestValidate(t *testing.T) {
	c, _ := New()
	require.NoError(t, c.Validate())

	c, _ = New()
	c.Log.Format = "wrongValue"
	require.Error(t, c.Validate())

	c, _ = New()
	c.Log.Level = "wrongValue"
	require.Error(t, c.Validate())
}
