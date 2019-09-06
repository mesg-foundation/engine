package cosmos

import (
	"encoding/json"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/mesg-foundation/engine/config"
	"github.com/mesg-foundation/engine/logger"
	tmconfig "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
	"github.com/tendermint/tendermint/types"
)

// NewNode creates a new Tendermint node from an App.
func NewNode(app *App, kb *Keybase, cfg *tmconfig.Config, ccfg *config.CosmosConfig) (*node.Node, error) {
	cdc := app.Cdc()

	// generate first user
	account, err := kb.GenerateAccount(ccfg.GenesisAccount.Name, ccfg.GenesisAccount.Mnemonic, ccfg.GenesisAccount.Password)
	if err != nil {
		return nil, err
	}

	// build a message to create validator
	msg := newMsgCreateValidator(
		sdktypes.ValAddress(account.GetAddress()),
		ed25519.PubKeyEd25519(ccfg.ValidatorPubKey),
		ccfg.GenesisAccount.Name,
	)
	signedTx, err := NewTxBuilder(cdc, 0, 0, kb, ccfg.ChainID).BuildAndSignStdTx(msg, ccfg.GenesisAccount.Name, ccfg.GenesisAccount.Password)
	if err != nil {
		return nil, err
	}

	// initialize app state with first validator
	appState, err := createAppState(app.DefaultGenesis(), cdc, account.GetAddress(), signedTx)
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

func createAppState(defaultGenesisŚtate map[string]json.RawMessage, cdc *codec.Codec, address sdktypes.AccAddress, signedStdTx authtypes.StdTx) (map[string]json.RawMessage, error) {
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

func newMsgCreateValidator(valAddr sdktypes.ValAddress, validatorPubKey ed25519.PubKeyEd25519, moniker string) sdktypes.Msg {
	return stakingtypes.NewMsgCreateValidator(
		valAddr,
		validatorPubKey,
		sdktypes.NewCoin(sdktypes.DefaultBondDenom, sdktypes.TokensFromConsensusPower(100)),
		stakingtypes.Description{
			Moniker: moniker,
			Details: "create-first-validator",
		},
		stakingtypes.NewCommissionRates(
			sdktypes.ZeroDec(),
			sdktypes.ZeroDec(),
			sdktypes.ZeroDec(),
		),
		sdktypes.NewInt(1),
	)
}
