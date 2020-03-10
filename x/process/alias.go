package process

import (
	"github.com/mesg-foundation/engine/x/process/internal/keeper"
	"github.com/mesg-foundation/engine/x/process/internal/types"
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

	QueryGetProcess    = types.QueryGetProcess
	QueryListProcesses = types.QueryListProcesses

	NewMsgCreateProcess = types.NewMsgCreateProcess
	NewMsgDeleteProcess = types.NewMsgDeleteProcess
)

// module types
type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
	Params       = types.Params

	MsgCreateProcess = types.MsgCreateProcess
	MsgDeleteProcess = types.MsgDeleteProcess
)
