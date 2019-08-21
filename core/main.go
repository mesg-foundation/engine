package main

import (
	"flag"
	"path/filepath"
	"strconv"
	"sync"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/config"
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/database"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/logger"
	"github.com/mesg-foundation/engine/scheduler"
	"github.com/mesg-foundation/engine/sdk"
	instancesdk "github.com/mesg-foundation/engine/sdk/instance"
	servicesdk "github.com/mesg-foundation/engine/sdk/service"
	"github.com/mesg-foundation/engine/server/grpc"
	"github.com/mesg-foundation/engine/store"
	"github.com/mesg-foundation/engine/tendermint"
	"github.com/mesg-foundation/engine/version"
	"github.com/mesg-foundation/engine/x/xerrors"
	"github.com/mesg-foundation/engine/x/xnet"
	"github.com/mesg-foundation/engine/x/xos"
	"github.com/mesg-foundation/engine/x/xsignal"
	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	db "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

var network = flag.Bool("experimental-network", false, "start the engine with the network")

func initialization(cfg *config.Config, network bool) (*sdk.SDK, *cosmos.App, error) {
	var app *cosmos.App

	// init container.
	c, err := container.New(cfg.Name)
	if err != nil {
		return nil, nil, err
	}

	var serviceSDK servicesdk.Service
	if network {
		// init cosmos app
		db, err := db.NewGoLevelDB("app", cfg.Cosmos.Path)
		if err != nil {
			return nil, nil, err
		}

		app = cosmos.New(logger.TendermintLogger(), db)

		// init cosmos stores
		paramsStoreKey := cosmostypes.NewKVStoreKey(params.StoreKey)
		app.RegisterStoreKey(paramsStoreKey)
		paramsTStoreKey := cosmostypes.NewTransientStoreKey(params.TStoreKey)
		app.RegisterTransientStoreKey(paramsTStoreKey)
		supplyStoreKey := cosmostypes.NewKVStoreKey(supply.StoreKey)
		app.RegisterStoreKey(supplyStoreKey)
		stakingStoreKey := cosmostypes.NewKVStoreKey(staking.StoreKey)
		app.RegisterStoreKey(stakingStoreKey)
		stakingTStoreKey := cosmostypes.NewTransientStoreKey(staking.TStoreKey)
		app.RegisterTransientStoreKey(stakingTStoreKey)
		distrStoreKey := cosmostypes.NewKVStoreKey(distribution.StoreKey)
		app.RegisterStoreKey(distrStoreKey)
		slashingStoreKey := cosmostypes.NewKVStoreKey(slashing.StoreKey)
		app.RegisterStoreKey(slashingStoreKey)

		// init cosmos keepers
		paramsKeeper := params.NewKeeper(
			app.Cdc(),
			paramsStoreKey,
			paramsTStoreKey,
			params.DefaultCodespace,
		)
		accountKeeper := auth.NewAccountKeeper(
			app.Cdc(),
			paramsStoreKey,
			paramsKeeper.Subspace(auth.DefaultParamspace),
			auth.ProtoBaseAccount,
		)
		bankKeeper := bank.NewBaseKeeper(
			accountKeeper,
			paramsKeeper.Subspace(bank.DefaultParamspace),
			bank.DefaultCodespace,
			nil,
		)
		supplyKeeper := supply.NewKeeper(
			app.Cdc(),
			supplyStoreKey,
			accountKeeper,
			bankKeeper,
			map[string][]string{
				auth.FeeCollectorName:     nil,
				distribution.ModuleName:   nil,
				staking.BondedPoolName:    {supply.Burner, supply.Staking},
				staking.NotBondedPoolName: {supply.Burner, supply.Staking},
			},
		)
		stakingKeeper := staking.NewKeeper(
			app.Cdc(),
			stakingStoreKey,
			stakingTStoreKey,
			supplyKeeper,
			paramsKeeper.Subspace(staking.DefaultParamspace),
			staking.DefaultCodespace,
		)
		distrKeeper := distribution.NewKeeper(
			app.Cdc(),
			distrStoreKey,
			paramsKeeper.Subspace(distribution.DefaultParamspace),
			&stakingKeeper,
			supplyKeeper,
			distribution.DefaultCodespace,
			auth.FeeCollectorName,
			nil,
		)
		slashingKeeper := slashing.NewKeeper(
			app.Cdc(),
			slashingStoreKey,
			&stakingKeeper,
			paramsKeeper.Subspace(slashing.DefaultParamspace),
			slashing.DefaultCodespace,
		)
		stakingKeeper = *stakingKeeper.SetHooks(
			staking.NewMultiStakingHooks(
				distrKeeper.Hooks(),
				slashingKeeper.Hooks()),
		)

		// init cosmos module
		app.RegisterModule(genaccounts.NewAppModule(accountKeeper))
		app.RegisterModule(genutil.NewAppModule(accountKeeper, stakingKeeper, app.DeliverTx))
		app.RegisterModule(auth.NewAppModule(accountKeeper))
		app.RegisterModule(bank.NewAppModule(bankKeeper, accountKeeper))
		app.RegisterModule(supply.NewAppModule(supplyKeeper, accountKeeper))
		app.RegisterModule(distribution.NewAppModule(distrKeeper, supplyKeeper))
		app.RegisterModule(slashing.NewAppModule(slashingKeeper, stakingKeeper))
		app.RegisterModule(staking.NewAppModule(stakingKeeper, distrKeeper, accountKeeper, supplyKeeper))

		app.RegisterOrderBeginBlocks(distribution.ModuleName, slashing.ModuleName)
		app.RegisterOrderEndBlocks(staking.ModuleName)

		app.RegisterOrderInitGenesis(
			genaccounts.ModuleName,
			distribution.ModuleName,
			staking.ModuleName,
			auth.ModuleName,
			bank.ModuleName,
			slashing.ModuleName,
			genutil.ModuleName,
		)

		// init service sdk
		serviceSDK = servicesdk.New(app, c)
	} else {
		serviceDB, err := leveldb.OpenFile(filepath.Join(cfg.Path, cfg.Database.ServiceRelativePath), nil)
		if err != nil {
			return nil, nil, err
		}
		serviceSDK = servicesdk.NewDeprecated(c, database.NewServiceKeeper(store.NewLevelDBStore(serviceDB)))
	}

	// init instance db.
	instanceDB, err := database.NewInstanceDB(filepath.Join(cfg.Path, cfg.Database.InstanceRelativePath))
	if err != nil {
		return nil, nil, err
	}

	// init execution db.
	executionDB, err := database.NewExecutionDB(filepath.Join(cfg.Path, cfg.Database.ExecutionRelativePath))
	if err != nil {
		return nil, nil, err
	}

	// init workflow db.
	workflowDB, err := database.NewWorkflowDB(filepath.Join(cfg.Path, cfg.Database.WorkflowRelativePath))
	if err != nil {
		return nil, nil, err
	}

	_, port, _ := xnet.SplitHostPort(cfg.Server.Address)

	// init sdk.
	return sdk.New(c, serviceSDK, instanceDB, executionDB, workflowDB, cfg.Name, strconv.Itoa(port)), app, nil
}

func deployCoreServices(cfg *config.Config, sdk *sdk.SDK) error {
	for _, serviceConfig := range cfg.SystemServices {
		logrus.WithField("module", "main").Infof("Deploying service %q", serviceConfig.Definition.Sid)
		srv, err := sdk.Service.Create(serviceConfig.Definition)
		if err != nil {
			existsError, ok := err.(*servicesdk.AlreadyExistsError)
			if ok {
				srv, err = sdk.Service.Get(existsError.Hash)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
		logrus.WithField("module", "main").Infof("Service %q deployed with hash %q", srv.Sid, srv.Hash)
		instance, err := sdk.Instance.Create(srv.Hash, xos.EnvMapToSlice(serviceConfig.Env))
		if err != nil {
			existsError, ok := err.(*instancesdk.AlreadyExistsError)
			if ok {
				instance, err = sdk.Instance.Get(existsError.Hash)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
		serviceConfig.Instance = instance
		logrus.WithField("module", "main").Infof("Instance started with hash %q", instance.Hash)
	}
	return nil
}

func stopRunningServices(sdk *sdk.SDK) error {
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

func main() {
	flag.Parse()
	cfg, err := config.Global()
	if err != nil {
		logrus.Fatalln(err)
	}

	sdk, app, err := initialization(cfg, *network)
	if err != nil {
		logrus.WithField("module", "main").Fatalln(err)
	}

	// init logger.
	logger.Init(cfg.Log.Format, cfg.Log.Level, cfg.Log.ForceColors)

	if *network {
		// load app
		err := app.Load()
		if err != nil {
			logrus.WithField("module", "main").Fatalln(err)
		}

		// create tendermint node
		node, err := tendermint.NewNode(app, cfg.Tendermint.Config, &cfg.Cosmos)
		if err != nil {
			logrus.WithField("module", "main").Fatalln(err)
		}

		// start tendermint node
		logrus.WithField("module", "main").WithField("seeds", cfg.Tendermint.P2P.Seeds).Info("starting tendermint node")
		if err := node.Start(); err != nil {
			logrus.WithField("module", "main").Fatalln(err)
		}
	}

	// init system services.
	// if err := deployCoreServices(cfg, sdk); err != nil {
	// 	logrus.WithField("module", "main").Fatalln(err)
	// }

	// init gRPC server.
	server := grpc.New(sdk)

	logrus.WithField("module", "main").Infof("starting MESG Engine version %s", version.Version)

	go func() {
		if err := server.Serve(cfg.Server.Address); err != nil {
			logrus.WithField("module", "main").Fatalln(err)
		}
	}()

	logrus.WithField("module", "main").Info("starting workflow engine")
	s := scheduler.New(sdk.Event, sdk.Execution, sdk.Workflow)
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
	server.Close()
}
