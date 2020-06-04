package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/x/credit/internal/types"
)

// Default parameter namespace
const (
	DefaultParamspace = types.ModuleName
)

// Minters are authorized account that can add credits to any address.
func (k Keeper) Minters(ctx sdk.Context) []sdk.AccAddress {
	var minters []sdk.AccAddress
	k.paramstore.Get(ctx, types.KeyMinters, &minters)
	return minters
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
