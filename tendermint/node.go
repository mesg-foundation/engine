package tendermint

import (
	"encoding/json"
	"time"

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
	tmclient "github.com/mesg-foundation/engine/tendermint/client"
	"github.com/sirupsen/logrus"
	tmconfig "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/types"
	db "github.com/tendermint/tm-db"
)

// NewNode retruns new tendermint node that runs the app.
func NewNode(cfg *tmconfig.Config, ccfg *config.CosmosConfig) (*node.Node, error) {
	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		return nil, err
	}

	me := privval.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile())

	// create user database and generate first user
	kb, err := NewFSKeybase(ccfg.Path)
	if err != nil {
		return nil, err
	}

	account, err := kb.GenerateAccount(ccfg.GenesisAccount.Name, ccfg.GenesisAccount.Mnemonic, ccfg.GenesisAccount.Password)
	if err != nil {
		return nil, err
	}

	// create app database and create an instance of the app
	db, err := db.NewGoLevelDB("app", ccfg.Path)
	if err != nil {
		return nil, err
	}

	logger := logger.TendermintLogger()
	cdc := app.MakeCodec()
	app := app.NewServiceApp(logger, db)

	// build a message to create validator
	msg := newMsgCreateValidator(
		sdktypes.ValAddress(account.GetAddress()),
		ed25519.PubKeyEd25519(ccfg.ValidatorPubKey),
		ccfg.GenesisAccount.Name,
	)
	logrus.WithField("msg", msg).Info("validator tx")

	// sign the message
	signedTx, err := signTransaction(
		cdc,
		kb,
		msg,
		ccfg.ChainID,
		ccfg.GenesisAccount.Name,
		ccfg.GenesisAccount.Password,
	)
	if err != nil {
		return nil, err
	}
	logrus.WithField("signedTx", signedTx).Info("signed tx")

	// initialize app state with first validator
	appState, err := createAppState(cdc, account.GetAddress(), signedTx)
	if err != nil {
		return nil, err
	}

	node, err := node.NewNode(cfg,
		me,
		nodeKey,
		proxy.NewLocalClientCreator(app),
		genesisLoader(cdc, appState, ccfg.ChainID, ccfg.GenesisTime),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(cfg.Instrumentation),
		logger,
	)
	if err != nil {
		return nil, err
	}

	// TODO: left only for tests
	go func() {
		client := tmclient.New(rpcclient.NewLocal(node), cdc)
		if services, err := client.ListServices(); err != nil {
			logrus.Error(err)
		} else {
			logrus.Info(services)
		}
	}()

	return node, nil
}

func signTransaction(cdc *codec.Codec, kb *Keybase, msg sdktypes.Msg, chainID, accountName, accountPassword string) (authtypes.StdTx, error) {
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
		chainID,
		"",
		sdktypes.NewCoins(),
		gasPrices,
	).WithKeybase(kb)

	return txBldr.SignStdTx(accountName, accountPassword, stdTx, false)
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
