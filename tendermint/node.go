package tendermint

import (
	"encoding/json"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authutils "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/mesg-foundation/engine/config"
	"github.com/mesg-foundation/engine/logger"
	"github.com/mesg-foundation/engine/tendermint/app"
	"github.com/sirupsen/logrus"
	tmconfig "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
	"github.com/tendermint/tendermint/types"
)

const (
	chainId         = "mesg-chain"
	accountName     = "bob"
	accountPassword = "12345678"
)

// NewNode retruns new tendermint node that runs the app.
func NewNode(cfg *tmconfig.Config, ccfg *config.CosmosConfig) (*node.Node, error) {
	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		return nil, err
	}

	me := privval.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile())

	// create user database and generate first user
	kb, err := NewFSKeybase(filepath.Join(ccfg.Path, "keybase"))
	if err != nil {
		return nil, err
	}

	account, err := kb.GenerateAccount(accountName, accountPassword)
	if err != nil {
		return nil, err
	}

	// create app database and create an instance of the app
	db, err := db.NewGoLevelDB("app", ccfg.Path)
	if err != nil {
		return nil, err
	}

	logger := logger.TendermintLogger()
	app, cdc := app.NewNameServiceApp(logger, db)

	// build a message to create validator
	msg := newMsgCreateValidator(sdktypes.ValAddress(account.GetAddress()), ed25519.PubKeyEd25519(ccfg.ValidatorPubKey))
	logrus.WithField("msg", msg).Info("validator tx")

	// sign the message
	fees := authtypes.NewStdFee(flags.DefaultGasLimit, sdktypes.NewCoins())
	gasPrices := sdktypes.DecCoins{}
	stdTx := authtypes.NewStdTx([]sdktypes.Msg{msg}, fees, []authtypes.StdSignature{}, "")

	txBldr := authtypes.NewTxBuilder(
		authutils.GetTxEncoder(cdc),
		0,
		0,
		flags.DefaultGasLimit,
		flags.DefaultGasAdjustment,
		true,
		chainId,
		"",
		sdktypes.NewCoins(),
		gasPrices,
	).WithKeybase(kb)

	signedTx, err := txBldr.SignStdTx(accountName, accountPassword, stdTx, false)
	if err != nil {
		return nil, err
	}
	logrus.WithField("signedTx", signedTx).Info("signed tx")

	// initialize app state with first validator
	appState, err := createAppState(cdc, account.GetAddress(), signedTx)
	if err != nil {
		return nil, err
	}

	return node.NewNode(cfg,
		me,
		nodeKey,
		proxy.NewLocalClientCreator(app),
		genesisLoader(cdc, appState),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(cfg.Instrumentation),
		logger,
	)
}

func createAppState(cdc *codec.Codec, address sdktypes.AccAddress, signedStdTx authtypes.StdTx) (map[string]json.RawMessage, error) {
	appState := app.ModuleBasics.DefaultGenesis()

	stakes := sdktypes.NewCoin(sdktypes.DefaultBondDenom, sdktypes.NewInt(100000000))
	genAcc := genaccounts.NewGenesisAccountRaw(address, sdktypes.NewCoins(stakes), sdktypes.NewCoins(), 0, 0, "", "")
	if err := genAcc.Validate(); err != nil {
		return nil, err
	}

	genstate, err := cdc.MarshalJSON(genaccounts.GenesisState([]genaccounts.GenesisAccount{genAcc}))
	if err != nil {
		return nil, err
	}
	appState[genaccounts.ModuleName] = genstate

	return genutil.SetGenTxsInAppGenesisState(cdc, appState, []authtypes.StdTx{signedStdTx})
}

func genesisLoader(cdc *codec.Codec, appState map[string]json.RawMessage) func() (*types.GenesisDoc, error) {
	return func() (*types.GenesisDoc, error) {
		appStateEncoded, err := cdc.MarshalJSON(appState)
		if err != nil {
			return nil, err
		}
		genesis := &types.GenesisDoc{
			ChainID:         chainId,
			ConsensusParams: types.DefaultConsensusParams(),
			AppState:        appStateEncoded,
		}
		if err := genesis.ValidateAndComplete(); err != nil {
			return nil, err
		}
		return genesis, nil
	}
}

func newMsgCreateValidator(valAddr sdktypes.ValAddress, validatorPubKey ed25519.PubKeyEd25519) sdktypes.Msg {
	return stakingtypes.NewMsgCreateValidator(
		valAddr,
		validatorPubKey,
		sdktypes.NewCoin(sdktypes.DefaultBondDenom, sdktypes.TokensFromConsensusPower(100)),
		stakingtypes.Description{
			Moniker: accountName,
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
