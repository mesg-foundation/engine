package ownership

import (
	"github.com/mesg-foundation/engine/x/ownership/internal/keeper"
	"github.com/mesg-foundation/engine/x/ownership/internal/types"
)

// const aliases
const (
	ModuleName   = types.ModuleName
	RouterKey    = types.RouterKey
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute
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
	QueryList = types.QueryList

	EventType                   = types.EventType
	AttributeKeyHash            = types.AttributeKeyHash
	AttributeKeyResourceHash    = types.AttributeKeyResourceHash
	AttributeKeyResourceType    = types.AttributeKeyResourceType
	AttributeKeyResourceAddress = types.AttributeKeyResourceAddress
	AttributeActionCreated      = types.AttributeActionCreated
	AttributeActionDeleted      = types.AttributeActionDeleted
)

// module types
type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState

	MsgWithdraw = types.MsgWithdraw
)
