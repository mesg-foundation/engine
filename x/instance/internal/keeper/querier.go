package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mesg-foundation/engine/cosmos/errors"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/x/instance/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier creates a new querier for instance clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryGetInstance:
			return getInstance(ctx, path[1:], k)
		case types.QueryListInstances:
			return listInstance(ctx, req, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown instance query endpoint")
		}
	}
}

func getInstance(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	if len(path) == 0 {
		return nil, errors.ErrMissingHash
	}
	hash, err := hash.Decode(path[0])
	if err != nil {
		return nil, sdkerrors.Wrap(errors.ErrValidation, err.Error())
	}

	instance, err := k.Get(ctx, hash)
	if err != nil {
		return nil, err
	}

	res, err := types.ModuleCdc.MarshalJSON(instance)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func listInstance(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var f *api.ListInstanceRequest
	if len(req.Data) > 0 {
		if err := types.ModuleCdc.UnmarshalJSON(req.Data, &f); err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}
	}

	instances, err := k.List(ctx, f.Filter)
	if err != nil {
		return nil, err
	}

	res, err := types.ModuleCdc.MarshalJSON(instances)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}
