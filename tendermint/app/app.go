package app

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/baseapp"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	servicesdk "github.com/mesg-foundation/engine/sdk/service"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

// TODO: sould be moved outside of app.gp
func BasicInit(moduleBasics ...module.AppModuleBasic) (module.BasicManager, *codec.Codec) {
	basicManager := module.NewBasicManager(
		append([]module.AppModuleBasic{
			genaccounts.AppModuleBasic{},
			genutil.AppModuleBasic{},
			auth.AppModuleBasic{},
			bank.AppModuleBasic{},
			staking.AppModuleBasic{},
			distr.AppModuleBasic{},
			params.AppModuleBasic{},
			slashing.AppModuleBasic{},
			supply.AppModuleBasic{},
		}, moduleBasics...)...,
	)

	cdc := codec.New()
	basicManager.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return basicManager, cdc
}

// NewServiceApp is a constructor function for ServiceApp
func NewServiceApp(cdc *codec.Codec, logger log.Logger, db dbm.DB, modules ...module.AppModule) *baseapp.BaseApp {

	// First define the top level codec that will be shared by the different modules
	// TODO: to delete
	// serviceSDK.SetCodec(cdc)

	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp("engine", logger, db, auth.DefaultTxDecoder(cdc))

	// Create store keys array
	storeKeys := []string{
		bam.MainStoreKey,
		auth.StoreKey,
		supply.StoreKey,
		staking.StoreKey,
		distr.StoreKey,
		params.StoreKey,
		slashing.StoreKey,
	}
	for _, module := range modules {
		storeKeys = append(storeKeys, module.Name())
	}

	keys := sdk.NewKVStoreKeys(storeKeys...)
	tkeys := sdk.NewTransientStoreKeys(staking.TStoreKey, params.TStoreKey)

	// Here you initialize your application with the store keys it requires

	// The ParamsKeeper handles parameter storage for the application
	paramsKeeper := params.NewKeeper(cdc, keys[params.StoreKey], tkeys[params.TStoreKey], params.DefaultCodespace)

	// The AccountKeeper handles address -> account lookups
	accountKeeper := auth.NewAccountKeeper(
		cdc,
		keys[params.StoreKey],
		paramsKeeper.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount,
	)

	// The BankKeeper allows you perform sdk.Coins interactions
	bankKeeper := bank.NewBaseKeeper(
		accountKeeper,
		paramsKeeper.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
		nil,
	)

	// The SupplyKeeper collects transaction fees and renders them to the fee distribution module
	supplyKeeper := supply.NewKeeper(
		cdc,
		keys[supply.StoreKey],
		accountKeeper,
		bankKeeper,
		map[string][]string{
			auth.FeeCollectorName:     nil,
			distr.ModuleName:          nil,
			staking.BondedPoolName:    {supply.Burner, supply.Staking},
			staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		},
	)

	// The staking keeper
	stakingKeeper := staking.NewKeeper(
		cdc,
		keys[staking.StoreKey],
		keys[staking.TStoreKey],
		supplyKeeper,
		paramsKeeper.Subspace(staking.DefaultParamspace),
		staking.DefaultCodespace,
	)

	distrKeeper := distr.NewKeeper(
		cdc,
		keys[distr.StoreKey],
		paramsKeeper.Subspace(distr.DefaultParamspace),
		&stakingKeeper,
		supplyKeeper,
		distr.DefaultCodespace,
		auth.FeeCollectorName,
		nil,
	)

	slashingKeeper := slashing.NewKeeper(
		cdc,
		keys[slashing.StoreKey],
		&stakingKeeper,
		paramsKeeper.Subspace(slashing.DefaultParamspace),
		slashing.DefaultCodespace,
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(
			distrKeeper.Hooks(),
			slashingKeeper.Hooks()),
	)

	// The serviceKeeper is the Keeper from the module for this tutorial
	// It handles interactions with the namestore
	serviceKeeper := servicesdk.NewKeeper(
		keys["service"],
		cdc,
	)

	mm := module.NewManager(
		genaccounts.NewAppModule(accountKeeper),
		genutil.NewAppModule(accountKeeper, stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(accountKeeper),
		bank.NewAppModule(bankKeeper, accountKeeper),
		//TODO: to finish: servicesdk.New(container, serviceKeeper),
		supply.NewAppModule(supplyKeeper, accountKeeper),
		distr.NewAppModule(distrKeeper, supplyKeeper),
		slashing.NewAppModule(slashingKeeper, stakingKeeper),
		staking.NewAppModule(stakingKeeper, distrKeeper, accountKeeper, supplyKeeper),
	)

	mm.SetOrderBeginBlockers(distr.ModuleName, slashing.ModuleName)
	mm.SetOrderEndBlockers(staking.ModuleName)

	// Sets the order of Genesis - Order matters, genutil is to always come last
	mm.SetOrderInitGenesis(
		genaccounts.ModuleName,
		distr.ModuleName,
		staking.ModuleName,
		auth.ModuleName,
		bank.ModuleName,
		slashing.ModuleName,
		serviceSDK.Name(),
		genutil.ModuleName,
	)

	// register all module routes and module queriers
	mm.RegisterRoutes(bApp.Router(), bApp.QueryRouter())

	// The initChainer handles translating the genesis.json file into initial state for the network
	bApp.SetInitChainer(func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		var genesisData map[string]json.RawMessage
		if err := cdc.UnmarshalJSON(req.AppStateBytes, &genesisData); err != nil {
			panic(err)
		}
		return mm.InitGenesis(ctx, genesisData)
	})
	bApp.SetBeginBlocker(func(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
		return mm.BeginBlock(ctx, req)
	})
	bApp.SetEndBlocker(func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		return mm.EndBlock(ctx, req)
	})

	// The AnteHandler handles signature verification and transaction pre-processing
	bApp.SetAnteHandler(
		auth.NewAnteHandler(
			accountKeeper,
			supplyKeeper,
			auth.DefaultSigVerificationGasConsumer,
		),
	)

	// initialize stores
	bApp.MountKVStores(keys)
	bApp.MountTransientStores(tkeys)

	if err := bApp.LoadLatestVersion(keys[bam.MainStoreKey]); err != nil {
		cmn.Exit(err.Error())
	}

	return bApp
}

// func (app *ServiceApp) LoadHeight(height int64) error {
// 	return app.LoadVersion(height, keys[bam.MainStoreKey])
// }
