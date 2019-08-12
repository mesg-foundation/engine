package tendermint

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authutils "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
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

var (
	chainId = "xxx"
)

// NewNode retruns new tendermint node that runs the app.
func NewNode(root, seeds, externalAddress string, validatorPukKey config.PubKeyEd25519) (*node.Node, error) {
	if err := os.MkdirAll(filepath.Join(root, "config"), 0755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Join(root, "data"), 0755); err != nil {
		return nil, err
	}

	cfg := tmconfig.DefaultConfig()
	cfg.P2P.PersistentPeers = seeds
	// cfg.P2P.Seeds = seeds
	// cfg.P2P.SeedMode = true
	cfg.P2P.AddrBookStrict = false
	cfg.P2P.AllowDuplicateIP = true
	cfg.P2P.ExternalAddress = externalAddress
	cfg.Consensus.TimeoutCommit = 10 * time.Second
	cfg.SetRoot(root)

	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		return nil, err
	}

	me := privval.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile())

	// init app
	db := db.NewMemDB()
	logger := logger.TendermintLogger()
	app := app.NewNameServiceApp(logger, db)

	appState, err := app.ExportInitialAppState()
	if err != nil {
		return nil, err
	}
	logrus.WithField("state", string(appState)).Info("state")

	// gen validator tx
	msg, err := newValidatorTx(validatorPukKey)
	if err != nil {
		return nil, err
	}
	logrus.WithField("msg", msg).Info("validator tx")
	// sign it
	fees := authtypes.NewStdFee(flags.DefaultGasLimit, sdktypes.NewCoins())
	gasPrices := sdktypes.NewDecCoins(sdktypes.NewCoins())
	stdTx := authtypes.NewStdTx([]sdktypes.Msg{msg}, fees, []authtypes.StdSignature{}, "")
	txBldr := authtypes.NewTxBuilder(
		authutils.GetTxEncoder(app.GetCodec()),
		0,
		0,
		flags.DefaultGasLimit,
		flags.DefaultGasAdjustment,
		true,
		chainId,
		"", // equal to "625b516112604044ea8e7d80e6c3544bd7766de8@10.0.2.15:26656" is genesis example
		sdktypes.NewCoins(),
		gasPrices,
	)
	signedTx, err := txBldr.SignStdTx("bob", "1", stdTx, true)
	if err != nil {
		return nil, err
	}
	logrus.WithField("signedTx", signedTx).Info("signed tx")

	// appMessage, err := genutil.GenAppStateFromConfig(cdc, config, initCfg, *genDoc, genAccIterator)
	// logrus.WithField("appMessage", appMessage).Info()

	return node.NewNode(cfg,
		me,
		nodeKey,
		proxy.NewLocalClientCreator(app),
		genesisLoader(appState),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(cfg.Instrumentation),
		logger,
	)
}

func genesisLoader(appState json.RawMessage) func() (*types.GenesisDoc, error) {
	return func() (*types.GenesisDoc, error) {
		genesis := &types.GenesisDoc{
			GenesisTime:     time.Date(2019, 8, 8, 0, 0, 0, 0, time.UTC),
			ChainID:         chainId,
			ConsensusParams: types.DefaultConsensusParams(),
			Validators:      []types.GenesisValidator{
				// 	{
				// 	Address: validator.Address(),
				// 	PubKey:  validator,
				// 	Power:   1,
				// 	Name:    "validator",
				// }
			},
			// AppState: []byte("{}"),
			AppState: appState,
		}
		if err := genesis.ValidateAndComplete(); err != nil {
			return nil, err
		}
		return genesis, nil
	}
}

func newValidatorTx(validatorPukKey config.PubKeyEd25519) (sdktypes.Msg, error) {
	emptyReturn := stakingtypes.MsgCreateValidator{}
	// default value see github.com/cosmos/cosmos-sdk@v0.36.0-rc1/x/staking/client/cli/tx.go

	validator := ed25519.PubKeyEd25519(validatorPukKey)
	// validatorAddress, err := sdktypes.ValAddressFromBech32("")

	defaultTokens := sdktypes.TokensFromConsensusPower(100)
	amount, err := sdktypes.ParseCoin(defaultTokens.String() + sdktypes.DefaultBondDenom)
	if err != nil {
		return emptyReturn, err
	}
	description := stakingtypes.NewDescription(
		"moniker-nico",
		"identify-nico",
		"website-nico",
		"details-nico",
	)

	rate, err := sdktypes.NewDecFromStr("0.1")
	if err != nil {
		return emptyReturn, err
	}
	maxRate, err := sdktypes.NewDecFromStr("0.2")
	if err != nil {
		return emptyReturn, err
	}
	maxChangeRate, err := sdktypes.NewDecFromStr("0.01")
	if err != nil {
		return emptyReturn, err
	}
	commissionRates := stakingtypes.NewCommissionRates(rate, maxRate, maxChangeRate)

	msg := stakingtypes.NewMsgCreateValidator(
		sdktypes.ValAddress(validator.Address()),
		validator,
		amount,
		description,
		commissionRates,
		sdktypes.NewInt(1),
	)
	return msg, nil
}
