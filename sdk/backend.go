package sdk

import (
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/cosmos"
	servicesdk "github.com/mesg-foundation/engine/sdk/service"
)

// Backend handles all the backend functions.
type Backend struct {
	Service *servicesdk.Backend
}

// NewBackend creates a new backend and init the sub-backend modules.
func NewBackend(appFactory *cosmos.AppFactory, c container.Container) *Backend {
	initDefaultCosmosModules(appFactory)
	service := servicesdk.NewBackend(appFactory, c)
	return &Backend{
		Service: service,
	}
}

func initDefaultCosmosModules(app *cosmos.AppFactory) {
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
}
