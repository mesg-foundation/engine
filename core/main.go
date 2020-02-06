package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"sync"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/gorilla/mux"
	"github.com/mesg-foundation/engine/app"
	"github.com/mesg-foundation/engine/codec"
	"github.com/mesg-foundation/engine/config"
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/ext/xerrors"
	"github.com/mesg-foundation/engine/ext/xnet"
	"github.com/mesg-foundation/engine/ext/xrand"
	"github.com/mesg-foundation/engine/ext/xsignal"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/logger"
	"github.com/mesg-foundation/engine/orchestrator"
	"github.com/mesg-foundation/engine/protobuf/api"
	enginesdk "github.com/mesg-foundation/engine/sdk"
	runnersdk "github.com/mesg-foundation/engine/sdk/runner"
	"github.com/mesg-foundation/engine/server/grpc"
	"github.com/mesg-foundation/engine/version"
	"github.com/sirupsen/logrus"
	rpcserver "github.com/tendermint/tendermint/rpc/lib/server"
	tmtypes "github.com/tendermint/tendermint/types"
	db "github.com/tendermint/tm-db"
)

func stopRunningServices(sdk *enginesdk.SDK, address string) error {
	runners, err := sdk.Runner.List(&runnersdk.Filter{Address: address})
	if err != nil {
		return err
	}
	var (
		runnersLen = len(runners)
		errC       = make(chan error, runnersLen)
		wg         sync.WaitGroup
	)
	wg.Add(runnersLen)
	for _, instance := range runners {
		go func(hash hash.Hash) {
			defer wg.Done()
			err := sdk.Runner.Delete(&api.DeleteRunnerRequest{
				Hash:       hash,
				DeleteData: false,
			})
			if err != nil {
				errC <- err
			}
		}(instance.Hash)
	}
	wg.Wait()
	close(errC)
	var errs xerrors.Errors
	for err := range errC {
		errs = append(errs, err)
	}
	return errs.ErrorOrNil()
}

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

func loadOrGenDevGenesis(kb *cosmos.Keybase, cfg *config.Config) (*tmtypes.GenesisDoc, error) {
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
	return cosmos.GenGenesis(kb, app.NewDefaultGenesisState(), cfg.DevGenesis.ChainID, cfg.DevGenesis.InitialBalances, cfg.DevGenesis.ValidatorDelegationCoin, cfg.Tendermint.Config.GenesisFile(), []cosmos.GenesisValidator{validator})
}

func main() {
	xrand.SeedInit()

	cfg, err := config.New()
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}
	cosmos.CustomizeConfig(cfg)

	// init logger.
	logger.Init(cfg.Log.Format, cfg.Log.Level, cfg.Log.ForceColors)

	// init basicManager
	// basicManager := enginesdk.NewBasicManager()

	// init tendermint logger
	tendermintLogger := logger.TendermintLogger()

	// init container.
	container, err := container.New(cfg.Name)
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	_, port, _ := xnet.SplitHostPort(cfg.Server.Address)

	// init app factory
	db, err := db.NewGoLevelDB("app", filepath.Join(cfg.Path, cfg.Cosmos.RelativePath))
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	initApp := app.NewInitApp(tendermintLogger, db, nil, true, 0, bam.SetMinGasPrices(cfg.Cosmos.MinGasPrices))

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
	genesis, err := loadOrGenDevGenesis(kb, cfg)
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	// create cosmos node
	node, err := cosmos.NewNode(initApp.BaseApp, cfg.Tendermint.Config, genesis)
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	// create cosmos client
	client, err := cosmos.NewClient(node, kb, genesis.ChainID, cfg.Account.Name, cfg.Account.Password, cfg.Cosmos.MinGasPrices)
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	// init sdk
	sdk := enginesdk.New(client, kb, container, cfg.Name, strconv.Itoa(port), cfg.IpfsEndpoint)

	// start tendermint node
	logrus.WithField("module", "main").WithField("seeds", cfg.Tendermint.Config.P2P.Seeds).Info("starting tendermint node")
	if err := node.Start(); err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	// init gRPC server.
	server := grpc.New(sdk, cfg, client)

	logrus.WithField("module", "main").Infof("starting MESG Engine version %s", version.Version)

	go func() {
		if err := server.Serve(cfg.Server.Address); err != nil {
			logrus.WithField("module", "main").Fatalln(err)
		}
	}()

	logrus.WithField("module", "main").Info("starting process engine")
	s := orchestrator.New(sdk.Event, sdk.Execution, sdk.Process, sdk.Runner)
	go func() {
		if err := s.Start(); err != nil {
			logrus.WithField("module", "main").Fatalln(err)
		}
	}()
	go func() {
		for err := range s.ErrC {
			logrus.WithField("module", "orchestrator").Warn(err)
		}
	}()

	logrus.WithField("module", "main").Info("starting lcd server")
	cfgLcd := rpcserver.DefaultConfig()
	lcdServer, err := rpcserver.Listen("tcp://[::]:1317", cfgLcd)
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	cliCtx := context.NewCLIContext().
		WithCodec(codec.Codec).
		WithClient(client).
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

	logrus.WithField("module", "main").Info("stopping lcd server")
	if err := lcdServer.Close(); err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	logrus.WithField("module", "main").Info("stopping running services")
	if err := stopRunningServices(sdk, acc.GetAddress().String()); err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	logrus.WithField("module", "main").Info("cleanup container")
	if err := container.Cleanup(); err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	logrus.WithField("module", "main").Info("stopping grpc server")
	server.Close()

	logrus.WithField("module", "main").Info("everything is stopped")
}
