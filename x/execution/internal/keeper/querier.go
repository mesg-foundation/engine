package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mesg-foundation/engine/cosmos/address"
	"github.com/mesg-foundation/engine/cosmos/errors"
	"github.com/mesg-foundation/engine/x/execution/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier creates a new querier for execution clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryGetExecution:
			return getExecution(ctx, k, path[1:])
		case types.QueryListExecution:
			return listExecution(ctx, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown execution query endpoint")
		}
	}
}

func getExecution(ctx sdk.Context, k Keeper, path []string) ([]byte, error) {
	if len(path) == 0 {
		return nil, errors.ErrMissingHash
	}
	hash, err := address.ExecAddressFromBech32(path[0])
	if err != nil {
		return nil, sdkerrors.Wrap(errors.ErrValidation, err.Error())
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

func listExecution(ctx sdk.Context, k Keeper) ([]byte, error) {
	es, err := k.List(ctx)
	if err != nil {
		return nil, err
	}

	res, err := types.ModuleCdc.MarshalJSON(es)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}
