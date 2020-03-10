package runner

import (
	"github.com/mesg-foundation/engine/x/runner/internal/keeper"
	"github.com/mesg-foundation/engine/x/runner/internal/types"
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

	ModuleCdc = types.ModuleCdc

	QueryGetRunner   = types.QueryGetRunner
	QueryListRunners = types.QueryListRunners

	NewMsgCreateRunner = types.NewMsgCreateRunner
	NewMsgDeleteRunner = types.NewMsgDeleteRunner
)

// module types
type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
	Params       = types.Params

	MsgCreateRunner = types.MsgCreateRunner
	MsgDeleteRunner = types.MsgDeleteRunner
)
