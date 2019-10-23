package cosmos

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/mesg-foundation/engine/logger"
	"github.com/stretchr/testify/require"
	db "github.com/tendermint/tm-db"
)

func appAndKbForTest() (*App, *Keybase, func(), error) {
	cosmosPath, _ := ioutil.TempDir("", "genesis-test-cosmos")
	closer := func() {
		os.RemoveAll(cosmosPath)
	}
	db, err := db.NewGoLevelDB("app", cosmosPath)
	if err != nil {
		return nil, nil, closer, err
	}
	appFactory := NewAppFactory(logger.TendermintLogger(), db)
	app, err := NewApp(appFactory)
	if err != nil {
		return nil, nil, closer, err
	}
	stakingtypes.RegisterCodec(app.Cdc())
	kb, err := NewKeybase(cosmosPath)
	if err != nil {
		return nil, nil, closer, err
	}
	return app, kb, closer, nil
}

func TestGenesis(t *testing.T) {
	app, kb, closer, err := appAndKbForTest()
	require.NoError(t, err)
	defer closer()
	path, _ := ioutil.TempDir("", "TestGenesis")
	defer os.RemoveAll(path)
	genesisPath := filepath.Join(path, "genesis.json")
	var (
		chainID                = "test-chainID"
		name                   = "name"
		password               = "pass"
		privValidatorKeyFile   = filepath.Join(path, "privValidatorKeyFile.json")
		privValidatorStateFile = filepath.Join(path, "privValidatorStateFile.json")
		nodeKeyFile            = filepath.Join(path, "nodeKeyFile.json")
		validators             = []GenesisValidator{}
	)
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
		genesis, err := GenGenesis(app, kb, chainID, genesisPath, validators)
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
