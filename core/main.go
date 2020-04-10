package main

import (
	"fmt"
	"path/filepath"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/gorilla/mux"
	"github.com/mesg-foundation/engine/app"
	"github.com/mesg-foundation/engine/config"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/event/publisher"
	"github.com/mesg-foundation/engine/ext/xrand"
	"github.com/mesg-foundation/engine/ext/xsignal"
	"github.com/mesg-foundation/engine/logger"
	"github.com/mesg-foundation/engine/orchestrator"
	"github.com/mesg-foundation/engine/server/grpc"
	"github.com/mesg-foundation/engine/version"
	"github.com/sirupsen/logrus"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	rpcserver "github.com/tendermint/tendermint/rpc/lib/server"
	tmtypes "github.com/tendermint/tendermint/types"
	db "github.com/tendermint/tm-db"
)

func loadOrGenConfigAccount(kb *cosmos.Keybase, cfg *config.Config) (keys.Info, error) {
	if cfg.Account.Mnemonic != "" {
		logrus.WithField("module", "main").Warn("Config account mnemonic presents. Generating account with it...")
		return kb.CreateAccount(cfg.Account.Name, cfg.Account.Mnemonic, "", cfg.Account.Password, keys.CreateHDPath(cfg.Account.Number, cfg.Account.Index).String(), cosmos.DefaultAlgo)
	}

	exist, err := kb.Exist(cfg.Account.Name)
	if err != nil {
		return nil, err
	}
	if exist {
		return kb.Get(cfg.Account.Name)
	}
	logrus.WithField("module", "main").Warn("Config account not found. Generating one for development...")
	mnemonic, err := kb.NewMnemonic()
	if err != nil {
		return nil, err
	}
	logrus.WithField("module", "main").WithFields(map[string]interface{}{
		"name":     cfg.Account.Name,
		"password": cfg.Account.Password,
		"mnemonic": mnemonic,
	}).Warn("Account")
	return kb.CreateAccount(cfg.Account.Name, mnemonic, "", cfg.Account.Password, keys.CreateHDPath(cfg.Account.Number, cfg.Account.Index).String(), cosmos.DefaultAlgo)
}

func loadOrGenDevGenesis(cdc *codec.Codec, kb *cosmos.Keybase, cfg *config.Config) (*tmtypes.GenesisDoc, error) {
	if cosmos.GenesisExist(cfg.Tendermint.Config.GenesisFile()) {
		return cosmos.LoadGenesis(cfg.Tendermint.Config.GenesisFile())
	}
	// generate dev genesis
	logrus.WithField("module", "main").Warn("Genesis file not found. Generating one for development...")
	validator, err := cosmos.NewGenesisValidator(kb,
		cfg.Account.Name,
		cfg.Account.Password,
		cfg.Tendermint.Config.PrivValidatorKeyFile(),
		cfg.Tendermint.Config.PrivValidatorStateFile(),
		cfg.Tendermint.Config.NodeKeyFile(),
	)
	if err != nil {
		return nil, err
	}
	logrus.WithField("module", "main").WithFields(map[string]interface{}{
		"nodeID": validator.NodeID,
		"peer":   fmt.Sprintf("%s@%s:26656", validator.NodeID, validator.Name),
	}).Warnln("Validator")
	return cosmos.GenGenesis(cdc, kb, app.NewDefaultGenesisState(), cfg.DevGenesis.ChainID, cfg.DevGenesis.InitialBalances, cfg.DevGenesis.ValidatorDelegationCoin, cfg.Tendermint.Config.GenesisFile(), []cosmos.GenesisValidator{validator})
}

//nolint:gocyclo
func main() {
	xrand.SeedInit()

	cfg, err := config.New()
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}
	cosmos.CustomizeConfig(cfg)

	// init logger.
	logger.Init(cfg.Log.Format, cfg.Log.Level, cfg.Log.ForceColors)

	// init tendermint logger
	tendermintLogger := logger.TendermintLogger()

	// init app factory
	db, err := db.NewGoLevelDB("app", filepath.Join(cfg.Path, cfg.Cosmos.RelativePath))
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	initApp, err := app.NewInitApp(tendermintLogger, db, nil, true, 0, bam.SetMinGasPrices(cfg.Cosmos.MinGasPrices))
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}
	cdc := initApp.Codec()

	// init key manager
	kb, err := cosmos.NewKeybase(filepath.Join(cfg.Path, cfg.Cosmos.RelativePath))
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	// gen config account
	acc, err := loadOrGenConfigAccount(kb, cfg)
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}
	logrus.WithField("address", acc.GetAddress().String()).Info("engine account")

	// load or gen genesis
	genesis, err := loadOrGenDevGenesis(cdc, kb, cfg)
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	// create cosmos node
	node, err := cosmos.NewNode(initApp.BaseApp, cfg.Tendermint.Config, genesis)
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}
	defer func() {
		logrus.WithField("module", "main").Info("stopping tendermint")
		if node.IsRunning() {
			if err := node.Stop(); err != nil {
				logrus.WithField("module", "main").Errorln(err)
			}
		}
	}()

	// create rpc client
	rpc, err := cosmos.NewRPC(rpcclient.NewLocal(node), cdc, kb, genesis.ChainID, cfg.Account.Name, cfg.Account.Password, cfg.Cosmos.MinGasPrices)
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	// init event publisher
	ep := publisher.New(rpc)

	// start tendermint node
	logrus.WithField("module", "main").WithField("seeds", cfg.Tendermint.Config.P2P.Seeds).Info("starting tendermint node")
	if err := node.Start(); err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	// init gRPC server.
	server := grpc.New(rpc, ep, cfg.AuthorizedPubKeys)
	logrus.WithField("module", "main").Infof("starting MESG Engine version %s", version.Version)
	defer func() {
		logrus.WithField("module", "main").Info("stopping grpc server")
		server.Close()
	}()

	go func() {
		if err := server.Serve(cfg.Server.Address); err != nil {
			logrus.WithField("module", "main").Fatalln(err)
		}
	}()

	logrus.WithField("module", "main").Info("starting process engine")
	orch := orchestrator.New(rpc, ep, cfg.DefaultExecutionPrice)
	defer func() {
		logrus.WithField("module", "main").Info("stopping orchestrator")
		orch.Stop()
	}()
	go func() {
		if err := orch.Start(); err != nil {
			logrus.WithField("module", "main").Fatalln(err)
		}
	}()
	go func() {
		for err := range orch.ErrC {
			logrus.WithField("module", "orchestrator").Warn(err)
		}
	}()

	logrus.WithField("module", "main").Info("starting lcd server")
	cfgLcd := rpcserver.DefaultConfig()
	lcdServer, err := rpcserver.Listen("tcp://[::]:1317", cfgLcd)
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}
	defer func() {
		logrus.WithField("module", "main").Info("stopping lcd server")
		if err := lcdServer.Close(); err != nil {
			logrus.WithField("module", "main").Errorln(err)
		}
	}()

	cliCtx := context.NewCLIContext().
		WithCodec(cdc).
		WithClient(rpc).
		WithTrustNode(true)
	mux := mux.NewRouter()
	cosmosclient.RegisterRoutes(cliCtx, mux)
	authrest.RegisterTxRoutes(cliCtx, mux)
	app.ModuleBasics.RegisterRESTRoutes(cliCtx, mux)
	go func() {
		if err := rpcserver.StartHTTPServer(lcdServer, mux, tendermintLogger, cfgLcd); err != nil {
			logrus.WithField("module", "main").Warnln(err) // not a fatal because closing the connection return an error here
		}
	}()

	<-xsignal.WaitForInterrupt()
}
