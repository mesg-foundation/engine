package serviceapp

import (
	"github.com/mesg-foundation/engine/tendermint/app/serviceapp/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey
)

var (
	NewMsgRemoveService = types.NewMsgRemoveService
	NewMsgSetService    = types.NewMsgSetService
	ModuleCdc           = types.ModuleCdc
	RegisterCodec       = types.RegisterCodec
)

type (
	MsgSetService        = types.MsgSetService
	MsgRemoveService     = types.MsgRemoveService
	QueryServiceResolve  = types.QueryServiceResolve
	QueryService         = types.QueryService
	QueryServicesResolve = types.QueryServicesResolve
	Service              = types.Service
)
