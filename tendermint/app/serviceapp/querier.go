package serviceapp

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	service := keeper.GetService(ctx, path[0])

	if service.Hash == "" {
		return []byte{}, sdk.ErrUnknownRequest("could not resolve service")
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, QueryServiceResolve{
		Service: QueryService{
			Hash:       path[0],
			Definition: service.Definition,
		},
	})
	if err != nil {
		return []byte{}, sdk.ErrInternal("could not unmarhsal service")
	}

	return res, nil
}

func queryServices(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	var services []Service
	for it := keeper.GetHashesIterator(ctx); it.Valid(); it.Next() {
		hash := string(it.Key())
		services = append(services, keeper.GetService(ctx, hash))
	}

	var qsr QueryServicesResolve
	for _, service := range services {
		qsr.Services = append(qsr.Services, QueryService{
			Hash:       service.Hash,
			Definition: service.Definition,
		})
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, qsr)
	if err != nil {
		return []byte{}, sdk.ErrInternal("could not unmarhsal services")
	}
	return res, nil
}
