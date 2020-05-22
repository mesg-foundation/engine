package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/x/execution/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier creates a new querier for execution clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryGet:
			return get(ctx, k, path[1:])
		case types.QueryList:
			return list(ctx, k, req.Data)
		case types.QueryParameters:
			return queryParameters(ctx, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown execution query endpoint")
		}
	}
}

func get(ctx sdk.Context, k Keeper, path []string) ([]byte, error) {
	if len(path) == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "missing hash")
	}
	hash, err := hash.Decode(path[0])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	e, err := k.Get(ctx, hash)
	if err != nil {
		return nil, err
	}

	res, err := types.ModuleCdc.MarshalJSON(e)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func list(ctx sdk.Context, k Keeper, data []byte) ([]byte, error) {
	var filter types.ListFilter
	if err := types.ModuleCdc.UnmarshalJSON(data, &filter); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	es, err := k.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	res, err := types.ModuleCdc.MarshalJSON(es)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryParameters(ctx sdk.Context, k Keeper) ([]byte, error) {
	res, err := types.ModuleCdc.MarshalJSON(k.GetParams(ctx))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}
