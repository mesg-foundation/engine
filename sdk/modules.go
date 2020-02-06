package sdk

import (
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp"
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/distribution"
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

// NewBasicManager returns a basic manager with all the app's modules.
func NewBasicManager() *module.BasicManager {
	basicManager := module.NewBasicManager(
		params.AppModuleBasic{},
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
	vesting.RegisterCodec(codec.Codec)
	basicManager.RegisterCodec(codec.Codec)
	return &basicManager
}

// NewApp returns a initialized and loaded cosmos app.
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
		supplyKeeper,
		paramsKeeper.Subspace(staking.DefaultParamspace),
	)
	distrKeeper := distribution.NewKeeper(
		codec.Codec,
		distrStoreKey,
		paramsKeeper.Subspace(distribution.DefaultParamspace),
		&stakingKeeper,
		supplyKeeper,
		auth.FeeCollectorName,
		modAccAddrs,
	)
	slashingKeeper := slashing.NewKeeper(
		codec.Codec,
		slashingStoreKey,
		&stakingKeeper,
		paramsKeeper.Subspace(slashing.DefaultParamspace),
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
	app.SetAppVersion(version.Version)

	// init module manager
	manager := module.NewManager(
		genutil.NewAppModule(accountKeeper, stakingKeeper, app.DeliverTx),
		auth.NewAppModule(accountKeeper),
		bank.NewAppModule(bankKeeper, accountKeeper),
		supply.NewAppModule(supplyKeeper, accountKeeper),
		distribution.NewAppModule(distrKeeper, accountKeeper, supplyKeeper, stakingKeeper),
		slashing.NewAppModule(slashingKeeper, accountKeeper, stakingKeeper),

		ownershipsdk.NewModule(ownershipKeeper),
		servicesdk.NewModule(serviceKeeper),
		instancesdk.NewModule(instanceKeeper),
		runnersdk.NewModule(runnerKeeper),
		executionsdk.NewModule(executionKeeper),
		processsdk.NewModule(processKeeper),

		staking.NewAppModule(stakingKeeper, accountKeeper, supplyKeeper),
	)
	manager.SetOrderBeginBlockers(distribution.ModuleName, slashing.ModuleName)
	manager.SetOrderEndBlockers(staking.ModuleName)
	manager.SetOrderInitGenesis(
		distribution.ModuleName,
		staking.ModuleName,
		auth.ModuleName,
		bank.ModuleName,
		slashing.ModuleName,

		// app module
		ownershipsdk.ModuleName,
		servicesdk.ModuleName,
		instancesdk.ModuleName,
		runnersdk.ModuleName,
		executionsdk.ModuleName,
		processsdk.ModuleName,

		// genutil should be last module
		supply.ModuleName,
		genutil.ModuleName,
	)

	// register app to manager
	manager.RegisterRoutes(app.Router(), app.QueryRouter())

	app.SetInitChainer(func(ctx cosmostypes.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		var genesisState simapp.GenesisState
		codec.Codec.MustUnmarshalJSON(req.AppStateBytes, &genesisState)
		return manager.InitGenesis(ctx, genesisState)
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
		processStoreKey,
	)

	app.SetAnteHandler(
		cosmostypes.ChainAnteDecorators(
			ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
			ante.NewMempoolFeeDecorator(),
			ante.NewValidateBasicDecorator(),
			ante.NewValidateMemoDecorator(accountKeeper),
			ante.NewConsumeGasForTxSizeDecorator(accountKeeper),
			ante.NewSetPubKeyDecorator(accountKeeper), // SetPubKeyDecorator must be called before all signature verification decorators
			ante.NewValidateSigCountDecorator(accountKeeper),
			ante.NewDeductFeeDecorator(accountKeeper, supplyKeeper),
			ante.NewSigGasConsumeDecorator(accountKeeper, auth.DefaultSigVerificationGasConsumer),
			ante.NewSigVerificationDecorator(accountKeeper),
			cosmos.NewForceCheckTxProxyDecorator(ante.NewIncrementSequenceDecorator(accountKeeper)), // innermost AnteDecorator
		),
		// Previous implementation:
		// auth.NewAnteHandler(accountKeeper, supplyKeeper, auth.DefaultSigVerificationGasConsumer)
	)
	if err := app.LoadLatestVersion(mainStoreKey); err != nil {
		return nil, err
	}

	return app, nil
}
