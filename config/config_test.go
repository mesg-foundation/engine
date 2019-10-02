package config

import (
	"os"
	"path/filepath"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/require"
)

func TestDefaultValue(t *testing.T) {
	home, _ := homedir.Dir()
	c, err := Default()
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

func TestLoad(t *testing.T) {
	snapsnot := map[string]string{
		"MESG_SERVER_ADDRESS":                 "",
		"MESG_LOG_FORMAT":                     "",
		"MESG_LOG_LEVEL":                      "",
		"MESG_LOG_FORCECOLORS":                "",
		"MESG_TENDERMINT_P2P_PERSISTENTPEERS": "",
		"MESG_COSMOS_GENESISVALIDATORTX":      "",
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
	os.Setenv("MESG_COSMOS_GENESISVALIDATORTX", `{"msg":[{"type":"cosmos-sdk/MsgCreateValidator","value":{"description":{"moniker":"bob","identity":"","website":"","details":"create-first-validator"},"commission":{"rate":"0.000000000000000000","max_rate":"0.000000000000000000","max_change_rate":"0.000000000000000000"},"min_self_delegation":"1","delegator_address":"cosmos1eav92wlgjjwkutl0s0u7nvdgc4m06zxtzdgwc0","validator_address":"cosmosvaloper1eav92wlgjjwkutl0s0u7nvdgc4m06zxt8eum5u","pubkey":"cosmosvalconspub1zcjduepqq0a87y3pur6vvzyp99t92me2zyxywz46kyq5vt7x2n4g987acmxszqey9p","value":{"denom":"stake","amount":"100000000"}}}],"fee":{"amount":[],"gas":"200000"},"signatures":[{"pub_key":{"type":"tendermint/PubKeySecp256k1","value":"A/p/EiHg9MYIgSlWVW8qEQxHCrqxAUYvxlTqgp/dxs3c"},"signature":"ehLZyLIARYDfYSolk9GtFb48g/5Mp33yHoIVuZZ2p78UkVABIKnFkfteaTqR0R06dTfbtcl2uTMhku0eT35Org=="}],"memo":"memo"}`)

	c, _ := Default()
	c.Load()
	require.Equal(t, "test_server_address", c.Server.Address)
	require.Equal(t, "test_log_format", c.Log.Format)
	require.Equal(t, "test_log_level", c.Log.Level)
	require.Equal(t, true, c.Log.ForceColors)
	require.Equal(t, "localhost", c.Tendermint.P2P.PersistentPeers)
	require.Equal(t, "memo", c.Cosmos.GenesisValidatorTx.Memo)
	require.Equal(t, uint64(200000), c.Cosmos.GenesisValidatorTx.Fee.Gas)
}

func TestValidate(t *testing.T) {
	c, _ := Default()
	c.Load()
	require.Error(t, c.Validate())
}

func TestStdTXDecodeNoSignersError(t *testing.T) {
	var tx StdTx
	require.Error(t, tx.Decode(`{"msg":[{"type":"cosmos-sdk/MsgCreateValidator","value":{"description":{"identity":"","website":"","details":"create-first-validator"},"pubkey":"cosmosvalconspub1zcjduepqq0a87y3pur6vvzyp99t92me2zyxywz46kyq5vt7x2n4g987acmxszqey9p","value":{"denom":"stake","amount":"100000000"}}}],"fee":{"amount":[],"gas":"200000"}}`))
}
