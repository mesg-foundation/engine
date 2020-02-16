package execution

import (
	"github.com/mesg-foundation/engine/x/execution/internal/keeper"
	"github.com/mesg-foundation/engine/x/execution/internal/types"
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

	QueryGetExecution  = types.QueryGetExecution
	QueryListExecution = types.QueryListExecution

	NewMsgCreateExecution = types.NewMsgCreateExecution
	NewMsgUpdateExecution = types.NewMsgUpdateExecution

	M = keeper.M
)

// module types
type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
	Params       = types.Params

	MsgCreateExecution = types.MsgCreateExecution
	MsgUpdateExecution = types.MsgUpdateExecution
)
