package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/x/execution/internal/types"
)

// Default parameter namespace
const (
	DefaultParamspace = types.ModuleName
)

// MinPrice - Minimum price of an execution
func (k Keeper) MinPrice(ctx sdk.Context) string {
	var coins string
	k.paramstore.Get(ctx, types.KeyMinPrice, &coins)
	return coins
}

// SetParams will populate all the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// GetParams returns all the params of the module
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	var params types.Params
	k.paramstore.GetParamSet(ctx, &params)
	return params
}
