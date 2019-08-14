package tendermint

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
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/logger"
	"github.com/mesg-foundation/engine/tendermint/app"
	servicetypes "github.com/mesg-foundation/engine/tendermint/app/serviceapp/types"
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
	// create user database and generate first user
	kb, err := NewKeybase(ccfg.Path)
	if err != nil {
		return nil, err
	}

	account, err := kb.GenerateAccount(ccfg.GenesisAccount.Name, ccfg.GenesisAccount.Mnemonic, ccfg.GenesisAccount.Password)
	if err != nil {
		return nil, err
	}

	cdc := app.MakeCodec()

	// build a message to create validator
	msg := newMsgCreateValidator(
		sdktypes.ValAddress(account.GetAddress()),
		ed25519.PubKeyEd25519(ccfg.ValidatorPubKey),
		ccfg.GenesisAccount.Name,
	)

	signedTx, err := NewTxBuilder(cdc, 0, 0, kb, ccfg.ChainID).
		DefaultSignStdTx(msg, ccfg.GenesisAccount.Name, ccfg.GenesisAccount.Password)
	if err != nil {
		return nil, err
	}

	// initialize app state with first validator
	appState, err := createAppState(cdc, account.GetAddress(), signedTx)
	if err != nil {
		return nil, err
	}

	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		return nil, err
	}

	// create app database and create an instance of the app
	db, err := db.NewGoLevelDB("app", ccfg.Path)
	if err != nil {
		return nil, err
	}

	app := app.NewServiceApp(logger.TendermintLogger(), db)

	node, err := node.NewNode(cfg,
		privval.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile()),
		nodeKey,
		proxy.NewLocalClientCreator(app),
		genesisLoader(cdc, appState, ccfg.ChainID, ccfg.GenesisTime),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(cfg.Instrumentation),
		app.Logger(),
	)
	if err != nil {
		return nil, err
	}

	// TODO: left only for tests
	go func() {
		client := tmclient.New(rpcclient.NewLocal(node), cdc)

		// add a service
		time.Sleep(22 * time.Second)
		msg := servicetypes.NewMsgSetService(hash.Int(2).String(), "{}", account.GetAddress())
		logrus.WithField("msg", msg).Warning("set service msg")

		accNumber, accSeq, err := authtypes.NewAccountRetriever(client).GetAccountNumberSequence(account.GetAddress())
		if err != nil {
			logrus.Error(err)
		}
		signedTx, err := NewTxBuilder(cdc, accNumber, accSeq, kb, ccfg.ChainID).
			BuildAndSign(ccfg.GenesisAccount.Name, ccfg.GenesisAccount.Password, []sdktypes.Msg{msg})
		if err != nil {
			logrus.Error(err)
		}

		logrus.WithField("signedTx", signedTx).Warning("set service signed tx")

		// broadcast the tx
		result, err := client.BroadcastTxSync(signedTx)
		if err != nil {
			logrus.Error(err)
		}
		logrus.WithField("result", result).Warning("tx broadcasted")

		msg = servicetypes.NewMsgSetService(hash.Int(1).String(), "{}", account.GetAddress())
		logrus.WithField("msg", msg).Warning("set service msg")

		signedTx, err = NewTxBuilder(cdc, accNumber, accSeq+1, kb, ccfg.ChainID).
			BuildAndSign(ccfg.GenesisAccount.Name, ccfg.GenesisAccount.Password, []sdktypes.Msg{msg})
		if err != nil {
			logrus.Error(err)
		}
		logrus.WithField("signedTx", signedTx).Warning("set service signed tx")

		// broadcast the tx
		result, err = client.BroadcastTxSync(signedTx)
		if err != nil {
			logrus.Error(err)
		}
		logrus.WithField("result", result).Warning("tx broadcasted")

		// fetch the service
		time.Sleep(11 * time.Second)
		if services, err := client.ListServices(); err != nil {
			logrus.Error(err)
		} else {
			logrus.Warning(services)
		}
	}()

	return node, nil
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

	// TODO: so we can use only one sign function
	// genesisState := genutil.GetGenesisStateFromAppState(cdc, appState)
	// genesisState.GenTxs = []json.RawMessage{}
	// return genutil.SetGenesisStateInAppState(cdc, appState, genesisState), nil

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
