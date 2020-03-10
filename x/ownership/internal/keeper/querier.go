package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mesg-foundation/engine/x/ownership/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier creates a new querier for ownership clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryListOwnerships:
			return listOwnerships(ctx, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown ownership query endpoint")
		}
	}
}

func listOwnerships(ctx sdk.Context, k Keeper) ([]byte, error) {
	ownerships, err := k.List(ctx)
	if err != nil {
		return nil, err
	}

	res, err := k.cdc.MarshalJSON(ownerships)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}
