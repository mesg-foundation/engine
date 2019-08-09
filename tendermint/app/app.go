package app

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

const appName = "mesg-app"

var (
	// ModuleBasics is BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager()

	mainStoreKey = sdk.NewKVStoreKey(baseapp.MainStoreKey)
)

// custom tx codec
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

// Extended ABCI application
type App struct {
	*bam.BaseApp
	cdc *codec.Codec

	// the module manager
	mm *module.Manager
}

// New returns a reference to an initialized App.
func New(logger log.Logger, db dbm.DB) (*App, *codec.Codec) {
	cdc := MakeCodec()

	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc))
	bApp.SetAppVersion(version.Version)

	var app = &App{
		BaseApp: bApp,
		cdc:     cdc,
	}

	app.mm = module.NewManager()

	app.MountStores(mainStoreKey)
	app.LoadLatestVersion(mainStoreKey)

	return app, cdc
}
