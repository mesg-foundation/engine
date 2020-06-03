package credit

import (
	"github.com/mesg-foundation/engine/x/credit/internal/keeper"
	"github.com/mesg-foundation/engine/x/credit/internal/types"
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

	ModuleCdc       = types.ModuleCdc
	QueryGet        = types.QueryGet
	QueryParameters = types.QueryParameters

	EventType                 = types.EventType
	AttributeActionAdded      = types.AttributeActionAdded
	AttributeActionSubtracted = types.AttributeActionSubtracted
)

// module types
type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
	Params       = types.Params

	MsgAdd = types.MsgAdd
)
