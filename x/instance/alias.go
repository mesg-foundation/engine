package instance

import (
	"github.com/mesg-foundation/engine/x/instance/internal/keeper"
	"github.com/mesg-foundation/engine/x/instance/internal/types"
)

// const aliases
const (
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	DefaultParamspace = types.DefaultParamspace
	QuerierRoute      = types.QuerierRoute
)

// functions and variable aliases
var (
	NewKeeper           = keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
	RegisterCodec       = types.RegisterCodec
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis

	ModuleCdc          = types.ModuleCdc
	QueryGetInstance   = types.QueryGetInstance
	QueryListInstances = types.QueryListInstances
)

// module types
type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
	Params       = types.Params
)
