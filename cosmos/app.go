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

// App is a Cosmos application that inherit from BaseApp.
type App struct {
	*baseapp.BaseApp
	cdc *codec.Codec

	modulesBasic       []module.AppModuleBasic
	modules            []module.AppModule
	storeKeys          map[string]*sdk.KVStoreKey
	transientStoreKeys map[string]*sdk.TransientStoreKey
	orderBeginBlockers []string
	orderEndBlockers   []string
	orderInitGenesis   []string
	anteHandler        sdk.AnteHandler
}

// NewApp initializes a new App.
func NewApp(logger log.Logger, db dbm.DB) *App {
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

// DefaultGenesis returns the default genesis from the basic manager.
func (a *App) DefaultGenesis() map[string]json.RawMessage {
	basicManager := module.NewBasicManager(a.modulesBasic...)
	basicManager.RegisterCodec(a.cdc)
	return basicManager.DefaultGenesis()
}

// RegisterModule registers a module to the app.
func (a *App) RegisterModule(module module.AppModule) {
	a.modulesBasic = append(a.modulesBasic, module)
	a.modules = append(a.modules, module)
}

// RegisterOrderInitGenesis sets the order of the modules when they are called to initialize the genesis.
func (a *App) RegisterOrderInitGenesis(moduleNames ...string) {
	a.orderInitGenesis = moduleNames
}

// RegisterOrderBeginBlocks sets the order of the modules when they are called on the begin block event.
func (a *App) RegisterOrderBeginBlocks(beginBlockers ...string) {
	a.orderBeginBlockers = beginBlockers
}

// RegisterOrderEndBlocks sets the order of the modules when they are called on the end block event.
func (a *App) RegisterOrderEndBlocks(endBlockers ...string) {
	a.orderEndBlockers = endBlockers
}

// RegisterStoreKey registers a store key to the app.
func (a *App) RegisterStoreKey(storeKey *sdk.KVStoreKey) {
	a.storeKeys[storeKey.Name()] = storeKey
}

// RegisterTransientStoreKey registers a transient store key to the app.
func (a *App) RegisterTransientStoreKey(transientStoreKey *sdk.TransientStoreKey) {
	a.transientStoreKeys[transientStoreKey.Name()] = transientStoreKey
}

// SetAnteHandler registers the authentication handler to the app.
func (a *App) SetAnteHandler(anteHandler sdk.AnteHandler) {
	a.anteHandler = anteHandler
}

// Cdc returns the codec of the app.
func (a *App) Cdc() *codec.Codec {
	return a.cdc
}

// Load creates the module manager, registers the modules to it, mounts the stores and finally load the app itself.
func (a *App) Load() error {
	// where all the magic happen
	// basically register everything on baseapp and load it

	mm := module.NewManager(a.modules...)
	mm.SetOrderBeginBlockers(a.orderBeginBlockers...)
	mm.SetOrderEndBlockers(a.orderEndBlockers...)

	// Sets the order of Genesis - Order matters, genutil is to always come last
	mm.SetOrderInitGenesis(a.orderInitGenesis...)

	// register all module routes and module queriers
	mm.RegisterRoutes(a.Router(), a.QueryRouter())

	// The initChainer handles translating the genesis.json file into initial state for the network
	a.SetInitChainer(func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		var genesisData map[string]json.RawMessage
		if err := a.cdc.UnmarshalJSON(req.AppStateBytes, &genesisData); err != nil {
			panic(err)
		}
		return mm.InitGenesis(ctx, genesisData)
	})
	a.SetBeginBlocker(mm.BeginBlock)
	a.SetEndBlocker(mm.EndBlock)

	// The AnteHandler handles signature verification and transaction pre-processing
	a.SetAnteHandler(a.anteHandler)

	// initialize stores
	a.MountKVStores(a.storeKeys)
	a.MountTransientStores(a.transientStoreKeys)

	return a.LoadLatestVersion(a.storeKeys[bam.MainStoreKey])
}
