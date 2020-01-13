package cosmos

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	cosmoscodec "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/mesg-foundation/engine/codec"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

// AppFactory is a Cosmos application factory.
type AppFactory struct {
	baseApp            *baseapp.BaseApp
	modulesBasic       []module.AppModuleBasic
	modules            []module.AppModule
	storeKeys          map[string]*sdk.KVStoreKey
	transientStoreKeys map[string]*sdk.TransientStoreKey
	orderBeginBlockers []string
	orderEndBlockers   []string
	orderInitGenesis   []string
	anteHandler        sdk.AnteHandler
}

// NewAppFactory returns a new AppFactory.
func NewAppFactory(logger log.Logger, db dbm.DB, minGasPrices string) *AppFactory {
	sdk.RegisterCodec(codec.Codec)
	cosmoscodec.RegisterCrypto(codec.Codec)

	return &AppFactory{
		baseApp: bam.NewBaseApp("engine", logger, db, auth.DefaultTxDecoder(codec.Codec), baseapp.SetMinGasPrices(minGasPrices)),
		modules: []module.AppModule{},
		storeKeys: map[string]*sdk.KVStoreKey{
			bam.MainStoreKey: sdk.NewKVStoreKey(bam.MainStoreKey),
		},
		transientStoreKeys: map[string]*sdk.TransientStoreKey{},
	}
}

// RegisterModule registers a module to the app.
func (a *AppFactory) RegisterModule(module module.AppModule) {
	a.modulesBasic = append(a.modulesBasic, module)
	a.modules = append(a.modules, module)
}

// RegisterOrderInitGenesis sets the order of the modules when they are called to initialize the genesis.
func (a *AppFactory) RegisterOrderInitGenesis(moduleNames ...string) {
	a.orderInitGenesis = moduleNames
}

// RegisterOrderBeginBlocks sets the order of the modules when they are called on the begin block event.
func (a *AppFactory) RegisterOrderBeginBlocks(beginBlockers ...string) {
	a.orderBeginBlockers = beginBlockers
}

// RegisterOrderEndBlocks sets the order of the modules when they are called on the end block event.
func (a *AppFactory) RegisterOrderEndBlocks(endBlockers ...string) {
	a.orderEndBlockers = endBlockers
}

// RegisterStoreKey registers a store key to the app.
func (a *AppFactory) RegisterStoreKey(storeKey *sdk.KVStoreKey) {
	a.storeKeys[storeKey.Name()] = storeKey
}

// RegisterTransientStoreKey registers a transient store key to the app.
func (a *AppFactory) RegisterTransientStoreKey(transientStoreKey *sdk.TransientStoreKey) {
	a.transientStoreKeys[transientStoreKey.Name()] = transientStoreKey
}

// SetAnteHandler registers the authentication handler to the app.
func (a *AppFactory) SetAnteHandler(anteHandler sdk.AnteHandler) {
	a.anteHandler = anteHandler
}

// DeliverTx implement baseApp.DeliverTx
func (a *AppFactory) DeliverTx(req abci.RequestDeliverTx) abci.ResponseDeliverTx {
	return a.baseApp.DeliverTx(req)
}
