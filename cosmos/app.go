package cosmos

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/baseapp"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	abci "github.com/tendermint/tendermint/abci/types"
)

// App is a loaded Cosmos application that inherit from BaseApp.
type App struct {
	*baseapp.BaseApp

	cdc          *codec.Codec
	basicManager module.BasicManager
}

// NewApp initializes a new App.
func NewApp(factory *AppFactory) (*App, error) {
	basicManager := module.NewBasicManager(factory.modulesBasic...)
	basicManager.RegisterCodec(factory.Cdc())

	a := &App{
		BaseApp:      factory.baseApp,
		basicManager: basicManager,
		cdc:          factory.Cdc(),
	}

	// Load creates the module manager, registers the modules to it, mounts the stores and finally load the app itself.
	mm := module.NewManager(factory.modules...)
	mm.SetOrderBeginBlockers(factory.orderBeginBlockers...)
	mm.SetOrderEndBlockers(factory.orderEndBlockers...)

	// Sets the order of Genesis - Order matters, genutil is to always come last
	mm.SetOrderInitGenesis(factory.orderInitGenesis...)

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
	a.SetAnteHandler(factory.anteHandler)

	// initialize stores
	a.MountKVStores(factory.storeKeys)
	a.MountTransientStores(factory.transientStoreKeys)

	if err := a.LoadLatestVersion(factory.storeKeys[bam.MainStoreKey]); err != nil {
		return nil, err
	}

	return a, nil
}

// DefaultGenesis returns the default genesis from the basic manager.
func (a *App) DefaultGenesis() map[string]json.RawMessage {
	return a.basicManager.DefaultGenesis()
}

// Cdc returns the codec of the app.
func (a *App) Cdc() *codec.Codec {
	return a.cdc
}
