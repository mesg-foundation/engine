package cosmos

import (
	"encoding/json"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/mesg-foundation/engine/config"
	"github.com/mesg-foundation/engine/logger"
	tmconfig "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
	"github.com/tendermint/tendermint/types"
)

// NewNode creates a new Tendermint node from an App.
func NewNode(app *App, cfg *tmconfig.Config, ccfg *config.CosmosConfig) (*node.Node, error) {
	cdc := app.Cdc()

	// initialize app state with first validator
	appState, err := createAppState(app.DefaultGenesis(), cdc, authtypes.StdTx(ccfg.GenesisValidatorTx))
	if err != nil {
		return nil, err
	}

	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		return nil, err
	}

	// init node
	return node.NewNode(
		cfg,
		privval.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile()),
		nodeKey,
		proxy.NewLocalClientCreator(app),
		genesisLoader(cdc, appState, ccfg.ChainID, ccfg.GenesisTime),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(cfg.Instrumentation),
		logger.TendermintLogger(),
	)
}

func createAppState(defaultGenesisŚtate map[string]json.RawMessage, cdc *codec.Codec, signedStdTx authtypes.StdTx) (map[string]json.RawMessage, error) {
	signers := signedStdTx.GetSigners()
	address := signers[len(signers)-1]
	stakes := sdktypes.NewCoin(sdktypes.DefaultBondDenom, sdktypes.NewInt(100000000))
	genAcc := genaccounts.NewGenesisAccountRaw(address, sdktypes.NewCoins(stakes), sdktypes.NewCoins(), 0, 0, "", "")
	if err := genAcc.Validate(); err != nil {
		return nil, err
	}

	genstate, err := cdc.MarshalJSON(genaccounts.GenesisState([]genaccounts.GenesisAccount{genAcc}))
	if err != nil {
		return nil, err
	}
	defaultGenesisŚtate[genaccounts.ModuleName] = genstate

	return genutil.SetGenTxsInAppGenesisState(cdc, defaultGenesisŚtate, []authtypes.StdTx{signedStdTx})
}

func genesisLoader(cdc *codec.Codec, appState map[string]json.RawMessage, chainID string, genesisTime time.Time) func() (*types.GenesisDoc, error) {
	return func() (*types.GenesisDoc, error) {
		appStateEncoded, err := cdc.MarshalJSON(appState)
		if err != nil {
			return nil, err
		}
		genesis := &types.GenesisDoc{
			GenesisTime:     genesisTime,
			ChainID:         chainID,
			ConsensusParams: types.DefaultConsensusParams(),
			AppState:        appStateEncoded,
		}
		if err := genesis.ValidateAndComplete(); err != nil {
			return nil, err
		}
		return genesis, nil
	}
}
