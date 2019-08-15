package serviceapp

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/hash"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the serviceapp Querier.
const (
	QueryServicePath  = "service"
	QueryServicesPath = "services"
)

// NewQuerier is the module level router for state queries.
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryServicePath:
			return queryService(ctx, path[1:], keeper)
		case QueryServicesPath:
			return queryServices(ctx, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown nameservice query endpoint")
		}
	}
}

func queryService(ctx sdk.Context, path []string, keeper Keeper) ([]byte, sdk.Error) {
	if len(path) == 0 {
		return nil, sdk.ErrUnknownRequest("no hash specified")
	}
	h, err := hash.Decode(path[0])
	if err != nil {
		return nil, sdk.ErrUnknownRequest(err.Error())
	}

	service := keeper.GetService(ctx, h)
	if service == nil {
		return nil, sdk.ErrUnknownRequest("service dosen't exist")
	}

	res, err := keeper.cdc.MarshalJSON(service)
	if err != nil {
		return nil, sdk.ErrInternal("could not unmarshal service")
	}
	return res, nil
}

func queryServices(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	services := keeper.GetServices(ctx)

	res, err := keeper.cdc.MarshalJSON(services)
	if err != nil {
		return nil, sdk.ErrInternal("could not unmarshal services")
	}
	return res, nil
}
