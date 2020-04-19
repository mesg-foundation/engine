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

	QueryGet  = types.QueryGet
	QueryList = types.QueryList

	M = keeper.M

	EventType                = types.EventType
	AttributeKeyHash         = types.AttributeKeyHash
	AttributeKeyAddress      = types.AttributeKeyAddress
	AttributeKeyExecutor     = types.AttributeKeyExecutor
	AttributeKeyProcess      = types.AttributeKeyProcess
	AttributeKeyInstance     = types.AttributeKeyInstance
	AttributeActionProposed  = types.AttributeActionProposed
	AttributeActionCreated   = types.AttributeActionCreated
	AttributeActionCompleted = types.AttributeActionCompleted
	AttributeActionFailed    = types.AttributeActionFailed
)

// module types
type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
	Params       = types.Params

	MsgCreate        = types.MsgCreate
	MsgUpdate        = types.MsgUpdate
	MsgUpdateOutputs = types.MsgUpdate_Outputs
	MsgUpdateError   = types.MsgUpdate_Error
)
