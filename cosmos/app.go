package cosmos

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/baseapp"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

// List of defaultbasic app module
// genaccounts.AppModuleBasic{},
// genutil.AppModuleBasic{},
// auth.AppModuleBasic{},
// bank.AppModuleBasic{},
// staking.AppModuleBasic{},
// distr.AppModuleBasic{},
// params.AppModuleBasic{},
// slashing.AppModuleBasic{},
// supply.AppModuleBasic{},

type App struct {
	// TODO: maybe it should inherit from baseapp in order to implement the right function for tendermint proxy client. let's try.

	logger log.Logger
	db     dbm.DB

	modulesBasic       []module.AppModuleBasic
	modules            []module.AppModule
	storeKeys          []string
	transientStoreKeys []string
	beginBlockers      []string
	endBlockers        []string
	modulesName        []string
	anteHandler        sdk.AnteHandler

	baseapp      *baseapp.BaseApp
	cdc          *codec.Codec
	basicManager module.BasicManager
}

func New(logger log.Logger, db dbm.DB) *App {
	cdc := codec.New()
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return &App{
		logger:  logger,
		db:      db,
		modules: []module.AppModule{},
		cdc:     cdc,
	}
}

func (a *App) BaseApp() *baseapp.BaseApp {
	return a.baseapp
}

func (a *App) DefaultGenesis() map[string]json.RawMessage {
	basicManager := module.NewBasicManager(a.modulesBasic...)
	basicManager.RegisterCodec(a.cdc)
	return basicManager.DefaultGenesis()
}

func (a *App) RegisterModule(module module.AppModule, storeKey, transientStoreKey string, isBeginBlocker, isEndBlocker bool) {
	moduleName := module.Name()
	a.modulesBasic = append(a.modulesBasic, module)
	a.modules = append(a.modules, module)
	a.modulesName = append(a.modulesName, moduleName)
	if storeKey != "" {
		a.storeKeys = append(a.storeKeys, storeKey)
	}
	if transientStoreKey != "" {
		a.transientStoreKeys = append(a.transientStoreKeys, transientStoreKey)
	}
	if isBeginBlocker {
		a.beginBlockers = append(a.beginBlockers, moduleName)
	}
	if isEndBlocker {
		a.endBlockers = append(a.endBlockers, moduleName)
	}
}

func (a *App) SetAnteHandler(anteHandler sdk.AnteHandler) {
	a.anteHandler = anteHandler
}

// TODO: is it really useful? better if not exported.
func (a *App) Cdc() *codec.Codec {
	return a.cdc
}

func (a *App) Load() {
	// where all the magic happen
	// basically register everything on baseapp and load it

	baseapp := bam.NewBaseApp("engine", a.logger, a.db, auth.DefaultTxDecoder(a.cdc))
	a.baseapp = baseapp

	mm := module.NewManager(a.modules...)
	mm.SetOrderBeginBlockers(a.beginBlockers...)
	mm.SetOrderEndBlockers(a.endBlockers...)

	// Sets the order of Genesis - Order matters, genutil is to always come last
	mm.SetOrderInitGenesis(a.modulesName...)

	// register all module routes and module queriers
	mm.RegisterRoutes(baseapp.Router(), baseapp.QueryRouter())

	// The initChainer handles translating the genesis.json file into initial state for the network
	baseapp.SetInitChainer(func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		var genesisData map[string]json.RawMessage
		if err := a.cdc.UnmarshalJSON(req.AppStateBytes, &genesisData); err != nil {
			panic(err)
		}
		return mm.InitGenesis(ctx, genesisData)
	})
	baseapp.SetBeginBlocker(func(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
		return mm.BeginBlock(ctx, req)
	})
	baseapp.SetEndBlocker(func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		return mm.EndBlock(ctx, req)
	})

	// The AnteHandler handles signature verification and transaction pre-processing
	baseapp.SetAnteHandler(a.anteHandler)

	// initialize stores
	storeKeys := sdk.NewKVStoreKeys(a.storeKeys...)
	baseapp.MountKVStores(storeKeys)
	transientStoreKeys := sdk.NewTransientStoreKeys(a.transientStoreKeys...)
	baseapp.MountTransientStores(transientStoreKeys)

	if err := baseapp.LoadLatestVersion(storeKeys[bam.MainStoreKey]); err != nil {
		cmn.Exit(err.Error())
	}
}
