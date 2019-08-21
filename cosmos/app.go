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
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

type App struct {
	*baseapp.BaseApp //TODO: revert if not useful

	logger log.Logger
	db     dbm.DB

	modulesBasic       []module.AppModuleBasic
	modules            []module.AppModule
	storeKeys          map[string]*sdk.KVStoreKey
	transientStoreKeys map[string]*sdk.TransientStoreKey
	orderBeginBlockers []string
	orderEndBlockers   []string
	orderInitGenesis   []string
	anteHandler        sdk.AnteHandler

	cdc          *codec.Codec
	basicManager module.BasicManager
}

func New(logger log.Logger, db dbm.DB) *App {
	cdc := codec.New()
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return &App{
		BaseApp: bam.NewBaseApp("engine", logger, db, auth.DefaultTxDecoder(cdc)),
		modules: []module.AppModule{},
		cdc:     cdc,
		storeKeys: map[string]*sdk.KVStoreKey{
			bam.MainStoreKey: sdk.NewKVStoreKey(bam.MainStoreKey),
		},
		transientStoreKeys: map[string]*sdk.TransientStoreKey{},
	}
}

func (a *App) DefaultGenesis() map[string]json.RawMessage {
	basicManager := module.NewBasicManager(a.modulesBasic...)
	basicManager.RegisterCodec(a.cdc)
	return basicManager.DefaultGenesis()
}

func (a *App) RegisterModule(module module.AppModule) {
	a.modulesBasic = append(a.modulesBasic, module)
	a.modules = append(a.modules, module)
}

func (a *App) RegisterOrderInitGenesis(moduleNames ...string) {
	a.orderInitGenesis = moduleNames
}

func (a *App) RegisterOrderBeginBlocks(beginBlockers ...string) {
	a.orderBeginBlockers = beginBlockers
}

func (a *App) RegisterOrderEndBlocks(endBlockers ...string) {
	a.orderEndBlockers = endBlockers
}

func (a *App) RegisterStoreKey(storeKey *sdk.KVStoreKey) {
	a.storeKeys[storeKey.Name()] = storeKey
}

func (a *App) RegisterTransientStoreKey(transientStoreKey *sdk.TransientStoreKey) {
	a.transientStoreKeys[transientStoreKey.Name()] = transientStoreKey
}

func (a *App) SetAnteHandler(anteHandler sdk.AnteHandler) {
	a.anteHandler = anteHandler
}

// TODO: is it really useful? better if not exported.
func (a *App) Cdc() *codec.Codec {
	return a.cdc
}

func (a *App) Load() error {
	// where all the magic happen
	// basically register everything on baseapp and load it

	mm := module.NewManager(a.modules...)
	mm.SetOrderBeginBlockers(a.orderBeginBlockers...)
	mm.SetOrderEndBlockers(a.orderEndBlockers...)

	// Sets the order of Genesis - Order matters, genutil is to always come last
	mm.SetOrderInitGenesis(a.orderInitGenesis...)

	// register all module routes and module queriers
	mm.RegisterRoutes(a.BaseApp.Router(), a.BaseApp.QueryRouter())

	// The initChainer handles translating the genesis.json file into initial state for the network
	a.BaseApp.SetInitChainer(func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		var genesisData map[string]json.RawMessage
		if err := a.cdc.UnmarshalJSON(req.AppStateBytes, &genesisData); err != nil {
			panic(err)
		}
		return mm.InitGenesis(ctx, genesisData)
	})
	a.BaseApp.SetBeginBlocker(func(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
		return mm.BeginBlock(ctx, req)
	})
	a.BaseApp.SetEndBlocker(func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		return mm.EndBlock(ctx, req)
	})

	// The AnteHandler handles signature verification and transaction pre-processing
	a.BaseApp.SetAnteHandler(a.anteHandler)

	// initialize stores
	a.BaseApp.MountKVStores(a.storeKeys)
	a.BaseApp.MountTransientStores(a.transientStoreKeys)

	return a.BaseApp.LoadLatestVersion(a.storeKeys[bam.MainStoreKey])
}
