package cosmos

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/baseapp"
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
)

type App struct {
	// TODO: maybe it should inherit from baseapp in order to implement the right function for tendermint proxy client. let's try.
	baseapp *baseapp.BaseApp

	modules []module.AppModule

	cdc          *codec.Codec
	basicManager module.BasicManager
}

func New(moduleBasics ...module.AppModuleBasic) *App {
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

	return &App{
		modules:      []module.AppModule{},
		cdc:          cdc,
		basicManager: basicManager,
	}
}

func (a *App) BaseApp() *baseapp.BaseApp {
	return a.baseapp
}

func (a *App) DefaultGenesis() map[string]json.RawMessage {
	// TODO: don't forget to take into account the module register with RegisterModule
	return a.basicManager.DefaultGenesis()
}

func (a *App) RegisterModule(module module.AppModule) {
	a.modules = append(a.modules, module)
}

// TODO: is it really useful? better if not exported.
func (a *App) Cdc() *codec.Codec {
	return a.cdc
}

func (a *App) Load() {
	// where all the magic happen
	// basically register everything on baseapp and load it
}
