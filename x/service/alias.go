package service

import (
	"github.com/mesg-foundation/engine/x/service/internal/keeper"
	"github.com/mesg-foundation/engine/x/service/internal/types"
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

	QueryGetService   = types.QueryGetService
	QueryListService  = types.QueryListService
	QueryHashService  = types.QueryHashService
	QueryExistService = types.QueryExistService

	NewMsgCreateService = types.NewMsgCreateService
)

// module types
type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
	Params       = types.Params

	MsgCreateService = types.MsgCreateService
)
