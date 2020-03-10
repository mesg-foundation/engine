package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mesg-foundation/engine/cosmos/address"
	"github.com/mesg-foundation/engine/cosmos/errors"
	"github.com/mesg-foundation/engine/x/process/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier creates a new querier for instance clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryGetProcess:
			return getProcess(ctx, path[1:], k)
		case types.QueryListProcesses:
			return listProcess(ctx, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown process query endpoint")
		}
	}
}

func getProcess(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	if len(path) == 0 {
		return nil, errors.ErrMissingHash
	}
	hash, err := address.ProcAddressFromBech32(path[0])
	if err != nil {
		return nil, sdkerrors.Wrap(errors.ErrValidation, err.Error())
	}

	instance, err := k.Get(ctx, hash)
	if err != nil {
		return nil, err
	}

	res, err := k.cdc.MarshalJSON(instance)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func listProcess(ctx sdk.Context, k Keeper) ([]byte, error) {
	instances, err := k.List(ctx)
	if err != nil {
		return nil, err
	}

	res, err := k.cdc.MarshalJSON(instances)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}
