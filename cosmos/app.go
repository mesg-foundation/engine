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
	baseapp *baseapp.BaseApp //TODO: revert if not useful
	cdc     *codec.Codec

	modulesBasic       []module.AppModuleBasic
	modules            []module.AppModule
	storeKeys          map[string]*sdk.KVStoreKey
	transientStoreKeys map[string]*sdk.TransientStoreKey
	orderBeginBlockers []string
	orderEndBlockers   []string
	orderInitGenesis   []string
	anteHandler        sdk.AnteHandler
}

func New(logger log.Logger, db dbm.DB) *App {
	cdc := codec.New()
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return &App{
		baseapp: bam.NewBaseApp("engine", logger, db, auth.DefaultTxDecoder(cdc)),
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
	mm.RegisterRoutes(a.baseapp.Router(), a.baseapp.QueryRouter())

	// The initChainer handles translating the genesis.json file into initial state for the network
	a.baseapp.SetInitChainer(func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		var genesisData map[string]json.RawMessage
		if err := a.cdc.UnmarshalJSON(req.AppStateBytes, &genesisData); err != nil {
			panic(err)
		}
		return mm.InitGenesis(ctx, genesisData)
	})
	a.baseapp.SetBeginBlocker(mm.BeginBlock)
	a.baseapp.SetEndBlocker(mm.EndBlock)

	// The AnteHandler handles signature verification and transaction pre-processing
	a.baseapp.SetAnteHandler(a.anteHandler)

	// initialize stores
	a.baseapp.MountKVStores(a.storeKeys)
	a.baseapp.MountTransientStores(a.transientStoreKeys)

	return a.baseapp.LoadLatestVersion(a.storeKeys[bam.MainStoreKey])
}

func (a *App) BaseApp() *baseapp.BaseApp {
	return a.baseapp
}

func (a *App) DeliverTx(req abci.RequestDeliverTx) (res abci.ResponseDeliverTx) {
	return a.baseapp.DeliverTx(req)
}
