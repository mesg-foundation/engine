package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/mesg-foundation/engine/config"
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/database"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/logger"
	"github.com/mesg-foundation/engine/orchestrator"
	enginesdk "github.com/mesg-foundation/engine/sdk"
	instancesdk "github.com/mesg-foundation/engine/sdk/instance"
	"github.com/mesg-foundation/engine/server/grpc"
	"github.com/mesg-foundation/engine/version"
	"github.com/mesg-foundation/engine/x/xerrors"
	"github.com/mesg-foundation/engine/x/xnet"
	"github.com/mesg-foundation/engine/x/xsignal"
	"github.com/sirupsen/logrus"
	tmtypes "github.com/tendermint/tendermint/types"
	db "github.com/tendermint/tm-db"
)

func initDatabases(cfg *config.Config) (*database.LevelDBInstanceDB, *database.LevelDBExecutionDB, *database.LevelDBProcessDB, error) {
	// init instance db.
	instanceDB, err := database.NewInstanceDB(filepath.Join(cfg.Path, cfg.Database.InstanceRelativePath))
	if err != nil {
		return nil, nil, nil, err
	}

	// init execution db.
	executionDB, err := database.NewExecutionDB(filepath.Join(cfg.Path, cfg.Database.ExecutionRelativePath))
	if err != nil {
		return nil, nil, nil, err
	}
	// init process db.
	processDB, err := database.NewProcessDB(filepath.Join(cfg.Path, cfg.Database.ProcessRelativePath))
	if err != nil {
		return nil, nil, nil, err
	}

	return instanceDB, executionDB, processDB, nil
}

func stopRunningServices(sdk *enginesdk.SDK) error {
	instances, err := sdk.Instance.List(&instancesdk.Filter{})
	if err != nil {
		return err
	}
	var (
		instancesLen = len(instances)
		errC         = make(chan error, instancesLen)
		wg           sync.WaitGroup
	)
	wg.Add(instancesLen)
	for _, instance := range instances {
		go func(hash hash.Hash) {
			defer wg.Done()
			err := sdk.Instance.Delete(hash, false)
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
	exist, err := kb.Exist(cfg.Account.Name)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, nil
	}
	logrus.WithField("module", "main").Warn("Config account not found. Generating one for development...")
	mnemonic, err := kb.NewMnemonic()
	if err != nil {
		return nil, err
	}
	return kb.CreateAccount(cfg.Account.Name, mnemonic, "", cfg.Account.Password, 0, 0)
}

func loadOrGenDevGenesis(app *cosmos.App, kb *cosmos.Keybase, cfg *config.Config) (*tmtypes.GenesisDoc, error) {
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
		"name":     validator.Name,
		"password": validator.Password,
		"nodeID":   validator.NodeID,
		"peer":     fmt.Sprintf("%s@%s:26656", validator.NodeID, validator.Name),
	}).Warnln("Dev validator")
	return cosmos.GenGenesis(app.Cdc(), kb, app.DefaultGenesis(), cfg.DevGenesis.ChainID, cfg.Tendermint.Config.GenesisFile(), []cosmos.GenesisValidator{validator})
}

func main() {
	cfg, err := config.New()
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	// init logger.
	logger.Init(cfg.Log.Format, cfg.Log.Level, cfg.Log.ForceColors)

	// init databases
	instanceDB, executionDB, processDB, err := initDatabases(cfg)
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	// init container.
	c, err := container.New(cfg.Name)
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	_, port, _ := xnet.SplitHostPort(cfg.Server.Address)

	// init app factory
	db, err := db.NewGoLevelDB("app", filepath.Join(cfg.Path, cfg.Cosmos.RelativePath))
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}
	// TODO: rename NewAppFactory to something else
	appFactory := cosmos.NewAppFactory(logger.TendermintLogger(), db)

	// register the backend modules to the app factory.
	// TODO: this is a mandatory call so it should return a new types required by cosmos.NewApp
	enginesdk.NewBackend(appFactory)

	// init cosmos app
	app, err := cosmos.NewApp(appFactory)
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	// init key manager
	kb, err := cosmos.NewKeybase(filepath.Join(cfg.Path, cfg.Cosmos.RelativePath))
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	// gen config account
	_, err = loadOrGenConfigAccount(kb, cfg)
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	// load or gen genesis
	genesis, err := loadOrGenDevGenesis(app, kb, cfg)
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	// create cosmos node
	node, err := cosmos.NewNode(app, cfg.Tendermint.Config, genesis)
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	// create cosmos client
	client := cosmos.NewClient(node, app.Cdc(), kb, genesis.ChainID)

	// init sdk
	sdk := enginesdk.New(client, app.Cdc(), kb, c, instanceDB, executionDB, processDB, cfg.Name, strconv.Itoa(port))

	// start tendermint node
	logrus.WithField("module", "main").WithField("seeds", cfg.Tendermint.Config.P2P.Seeds).Info("starting tendermint node")
	if err := node.Start(); err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	// init gRPC server.
	server := grpc.New(sdk, cfg)

	logrus.WithField("module", "main").Infof("starting MESG Engine version %s", version.Version)

	go func() {
		if err := server.Serve(cfg.Server.Address); err != nil {
			logrus.WithField("module", "main").Fatalln(err)
		}
	}()

	logrus.WithField("module", "main").Info("starting process engine")
	s := orchestrator.New(sdk.Event, sdk.Execution, sdk.Process)
	go func() {
		if err := s.Start(); err != nil {
			logrus.WithField("module", "main").Fatalln(err)
		}
	}()
	go func() {
		for err := range s.ErrC {
			logrus.WithField("module", "main").Warn(err)
		}
	}()

	<-xsignal.WaitForInterrupt()
	if err := stopRunningServices(sdk); err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	if err := c.Cleanup(); err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}
	server.Close()
}
