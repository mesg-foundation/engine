package sdk

import (
	"encoding/json"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/mesg-foundation/engine/codec"
	"github.com/mesg-foundation/engine/cosmos"
	executionsdk "github.com/mesg-foundation/engine/sdk/execution"
	instancesdk "github.com/mesg-foundation/engine/sdk/instance"
	ownershipsdk "github.com/mesg-foundation/engine/sdk/ownership"
	processsdk "github.com/mesg-foundation/engine/sdk/process"
	runnersdk "github.com/mesg-foundation/engine/sdk/runner"
	servicesdk "github.com/mesg-foundation/engine/sdk/service"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

func NewBasicManager() *module.BasicManager {
	basicManager := module.NewBasicManager(
		params.AppModuleBasic{},
		genaccounts.AppModuleBasic{},
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		supply.AppModuleBasic{},
		distribution.AppModuleBasic{},
		slashing.AppModuleBasic{},
		staking.AppModuleBasic{},

		cosmos.NewAppModuleBasic(ownershipsdk.ModuleName),
		cosmos.NewAppModuleBasic(servicesdk.ModuleName),
		cosmos.NewAppModuleBasic(instancesdk.ModuleName),
		cosmos.NewAppModuleBasic(runnersdk.ModuleName),
		cosmos.NewAppModuleBasic(executionsdk.ModuleName),
		cosmos.NewAppModuleBasic(processsdk.ModuleName),
	)
	basicManager.RegisterCodec(codec.Codec)
	return &basicManager
}

func NewApp(logger log.Logger, db dbm.DB, minGasPrices string) (*bam.BaseApp, error) {
	// init cosmos stores
	mainStoreKey := cosmostypes.NewKVStoreKey(bam.MainStoreKey)
	paramsStoreKey := cosmostypes.NewKVStoreKey(params.StoreKey)
	paramsTStoreKey := cosmostypes.NewTransientStoreKey(params.TStoreKey)
	supplyStoreKey := cosmostypes.NewKVStoreKey(supply.StoreKey)
	stakingStoreKey := cosmostypes.NewKVStoreKey(staking.StoreKey)
	stakingTStoreKey := cosmostypes.NewTransientStoreKey(staking.TStoreKey)
	distrStoreKey := cosmostypes.NewKVStoreKey(distribution.StoreKey)
	slashingStoreKey := cosmostypes.NewKVStoreKey(slashing.StoreKey)

	ownershipStoreKey := cosmostypes.NewKVStoreKey(ownershipsdk.ModuleName)
	serviceStoreKey := cosmostypes.NewKVStoreKey(servicesdk.ModuleName)
	instanceStoreKey := cosmostypes.NewKVStoreKey(instancesdk.ModuleName)
	runnerStoreKey := cosmostypes.NewKVStoreKey(runnersdk.ModuleName)
	executionStoreKey := cosmostypes.NewKVStoreKey(executionsdk.ModuleName)
	processStoreKey := cosmostypes.NewKVStoreKey(processsdk.ModuleName)

	// account permissions
	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		distribution.ModuleName:   nil,
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
	}
	// Module Accounts Addresses
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	// init cosmos keepers
	paramsKeeper := params.NewKeeper(
		codec.Codec,
		paramsStoreKey,
		paramsTStoreKey,
		params.DefaultCodespace,
	)
	accountKeeper := auth.NewAccountKeeper(
		codec.Codec,
		paramsStoreKey,
		paramsKeeper.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount,
	)
	bankKeeper := bank.NewBaseKeeper(
		accountKeeper,
		paramsKeeper.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
		modAccAddrs,
	)
	supplyKeeper := supply.NewKeeper(
		codec.Codec,
		supplyStoreKey,
		accountKeeper,
		bankKeeper,
		maccPerms,
	)
	stakingKeeper := staking.NewKeeper(
		codec.Codec,
		stakingStoreKey,
		stakingTStoreKey,
		supplyKeeper,
		paramsKeeper.Subspace(staking.DefaultParamspace),
		staking.DefaultCodespace,
	)
	distrKeeper := distribution.NewKeeper(
		codec.Codec,
		distrStoreKey,
		paramsKeeper.Subspace(distribution.DefaultParamspace),
		&stakingKeeper,
		supplyKeeper,
		distribution.DefaultCodespace,
		auth.FeeCollectorName,
		modAccAddrs,
	)
	slashingKeeper := slashing.NewKeeper(
		codec.Codec,
		slashingStoreKey,
		&stakingKeeper,
		paramsKeeper.Subspace(slashing.DefaultParamspace),
		slashing.DefaultCodespace,
	)
	stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(
			distrKeeper.Hooks(),
			slashingKeeper.Hooks(),
		),
	)
	ownershipKeeper := ownershipsdk.NewKeeper(ownershipStoreKey)
	serviceKeeper := servicesdk.NewKeeper(serviceStoreKey, ownershipKeeper)
	instanceKeeper := instancesdk.NewKeeper(instanceStoreKey)
	runnerKeeper := runnersdk.NewKeeper(runnerStoreKey, instanceKeeper)
	processKeeper := processsdk.NewKeeper(processStoreKey, ownershipKeeper, instanceKeeper)
	executionKeeper := executionsdk.NewKeeper(executionStoreKey, serviceKeeper, instanceKeeper, runnerKeeper, processKeeper)

	// init app
	// TODO: engine should be in config
	app := bam.NewBaseApp("engine", logger, db, auth.DefaultTxDecoder(codec.Codec), bam.SetMinGasPrices(minGasPrices))

	// init module manager
	manager := module.NewManager(
		genaccounts.NewAppModule(accountKeeper),
		genutil.NewAppModule(accountKeeper, stakingKeeper, app.DeliverTx),
		auth.NewAppModule(accountKeeper),
		bank.NewAppModule(bankKeeper, accountKeeper),
		supply.NewAppModule(supplyKeeper, accountKeeper),
		distribution.NewAppModule(distrKeeper, supplyKeeper),
		slashing.NewAppModule(slashingKeeper, stakingKeeper),
		staking.NewAppModule(stakingKeeper, distrKeeper, accountKeeper, supplyKeeper),

		ownershipsdk.NewModule(ownershipKeeper),
		servicesdk.NewModule(serviceKeeper),
		instancesdk.NewModule(instanceKeeper),
		runnersdk.NewModule(runnerKeeper),
		executionsdk.NewModule(executionKeeper),
		processsdk.NewModule(processKeeper),
	)
	manager.SetOrderBeginBlockers(distribution.ModuleName, slashing.ModuleName)
	manager.SetOrderEndBlockers(staking.ModuleName)
	manager.SetOrderInitGenesis(
		genaccounts.ModuleName,
		distribution.ModuleName,
		staking.ModuleName,
		auth.ModuleName,
		bank.ModuleName,
		slashing.ModuleName,
		supply.ModuleName,

		// app module
		ownershipsdk.ModuleName,
		servicesdk.ModuleName,
		instancesdk.ModuleName,
		runnersdk.ModuleName,
		executionsdk.ModuleName,
		processsdk.ModuleName,

		// genutil should be last module
		genutil.ModuleName,
	)

	// register app to manager
	manager.RegisterRoutes(app.Router(), app.QueryRouter())
	app.SetInitChainer(func(ctx cosmostypes.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		var genesisData map[string]json.RawMessage
		if err := codec.UnmarshalJSON(req.AppStateBytes, &genesisData); err != nil {
			panic(err)
		}
		return manager.InitGenesis(ctx, genesisData)
	})
	app.SetBeginBlocker(manager.BeginBlock)
	app.SetEndBlocker(manager.EndBlock)

	app.MountStores(
		mainStoreKey,
		paramsStoreKey,
		paramsTStoreKey,
		supplyStoreKey,
		stakingStoreKey,
		stakingTStoreKey,
		distrStoreKey,
		slashingStoreKey,

		ownershipStoreKey,
		serviceStoreKey,
		instanceStoreKey,
		runnerStoreKey,
		executionStoreKey,
	)
	app.SetAnteHandler(auth.NewAnteHandler(accountKeeper, supplyKeeper, auth.DefaultSigVerificationGasConsumer))
	if err := app.LoadLatestVersion(mainStoreKey); err != nil {
		return nil, err
	}

	return app, nil
}
