package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/x/service/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier creates a new querier for service clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryGet:
			return get(ctx, k, path[1:])
		case types.QueryList:
			return list(ctx, k)
		case types.QueryExist:
			return exist(ctx, k, path[1:])
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown service query endpoint")
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

	srv, err := k.Get(ctx, hash)
	if err != nil {
		return nil, err
	}

	res, err := types.ModuleCdc.MarshalJSON(srv)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func list(ctx sdk.Context, k Keeper) ([]byte, error) {
	srvs, err := k.List(ctx)
	if err != nil {
		return nil, err
	}

	res, err := types.ModuleCdc.MarshalJSON(srvs)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func exist(ctx sdk.Context, k Keeper, path []string) ([]byte, error) {
	if len(path) == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "missing hash")
	}
	hash, err := hash.Decode(path[0])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	exists, err := k.Exists(ctx, hash)
	if err != nil {
		return nil, err
	}

	res, err := types.ModuleCdc.MarshalJSON(exists)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}
