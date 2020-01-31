package cosmos

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/mesg-foundation/engine/codec"
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
	stakingtypes.RegisterCodec(codec.Codec)
	// variables
	var (
		chainID                 = "test-chainID"
		initialBalances         = "100amesg"
		validatorDelegationCoin = "1amesg"
		name                    = "name"
		password                = "pass"
		privValidatorKeyFile    = filepath.Join(path, "privValidatorKeyFile.json")
		privValidatorStateFile  = filepath.Join(path, "privValidatorStateFile.json")
		nodeKeyFile             = filepath.Join(path, "nodeKeyFile.json")
		genesisPath             = filepath.Join(path, "genesis.json")
		validators              = []GenesisValidator{}
		defaultGenesisState     = map[string]json.RawMessage{}
	)
	// init account
	mnemonic, _ := kb.NewMnemonic()
	kb.CreateAccount(name, mnemonic, "", password, 0, 0)
	// start tests
	t.Run("generate validator", func(t *testing.T) {
		v, err := NewGenesisValidator(kb, name, password, privValidatorKeyFile, privValidatorStateFile, nodeKeyFile)
		validators = append(validators, v)
		require.NoError(t, err)
		require.Equal(t, name, v.Name)
		require.Equal(t, password, v.Password)
		require.NotEmpty(t, v.ValPubKey)
		require.NotEmpty(t, v.NodeID)
		require.FileExists(t, privValidatorKeyFile)
		require.FileExists(t, privValidatorStateFile)
		require.FileExists(t, nodeKeyFile)
	})
	t.Run("genesis doesn't exist", func(t *testing.T) {
		require.False(t, GenesisExist(genesisPath))
	})
	t.Run("generate genesis", func(t *testing.T) {
		genesis, err := GenGenesis(kb, defaultGenesisState, chainID, initialBalances, validatorDelegationCoin, genesisPath, validators)
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
