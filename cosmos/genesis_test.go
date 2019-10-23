package cosmos

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	// path
	path, _ := ioutil.TempDir("", "TestGenesis")
	defer os.RemoveAll(path)
	// keybase
	kb, err := NewKeybase(filepath.Join(path, "kb"))
	require.NoError(t, err)
	// codec
	cdc := codec.New()
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	stakingtypes.RegisterCodec(cdc)
	// variables
	var (
		chainID                = "test-chainID"
		name                   = "name"
		password               = "pass"
		privValidatorKeyFile   = filepath.Join(path, "privValidatorKeyFile.json")
		privValidatorStateFile = filepath.Join(path, "privValidatorStateFile.json")
		nodeKeyFile            = filepath.Join(path, "nodeKeyFile.json")
		genesisPath            = filepath.Join(path, "genesis.json")
		validators             = []GenesisValidator{}
		defaultGenesisState    = map[string]json.RawMessage{}
	)
	// start tests
	t.Run("generate validator", func(t *testing.T) {
		v, err := NewGenesisValidator(kb, name, password, privValidatorKeyFile, privValidatorStateFile, nodeKeyFile)
		validators = append(validators, v)
		require.NoError(t, err)
		require.Equal(t, name, v.Name)
		require.Equal(t, password, v.Password)
		require.NotEmpty(t, v.Address)
		require.NotEmpty(t, v.Mnemonic)
		require.NotEmpty(t, v.ValPubKey)
		require.NotEmpty(t, v.NodeID)
		require.FileExists(t, privValidatorKeyFile)
		require.FileExists(t, privValidatorStateFile)
		require.FileExists(t, nodeKeyFile)
		acc, err := kb.GetByAddress(v.Address)
		require.NoError(t, err)
		require.Equal(t, name, acc.GetName())
	})
	t.Run("genesis doesn't exist", func(t *testing.T) {
		require.False(t, GenesisExist(genesisPath))
	})
	t.Run("generate genesis", func(t *testing.T) {
		genesis, err := GenGenesis(cdc, kb, defaultGenesisState, chainID, genesisPath, validators)
		require.NoError(t, err)
		require.NotEmpty(t, genesis)
	})
	t.Run("load genesis", func(t *testing.T) {
		genesis, err := LoadGenesis(genesisPath)
		require.NoError(t, err)
		require.NotEmpty(t, genesis)
		require.Equal(t, chainID, genesis.ChainID)
	})
	t.Run("genesis exist", func(t *testing.T) {
		require.True(t, GenesisExist(genesisPath))
	})
}
