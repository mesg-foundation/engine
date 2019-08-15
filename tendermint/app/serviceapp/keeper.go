package serviceapp

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/service"
)

// Keeper maintains the link to data storage and exposes
// getter/setter methods for the various parts of the service state machine.
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
}

// NewKeeper creates new instances of the service Keeper.
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

// Gets the entire Whois metadata struct for a name
func (k Keeper) GetService(ctx sdk.Context, hash hash.Hash) *service.Service {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(hash) {
		return nil
	}

	b := store.Get(hash)
	var service service.Service
	k.cdc.MustUnmarshalBinaryBare(b, &service)
	return &service
}

// Sets the entire Whois metadata struct for a name
func (k Keeper) SetService(ctx sdk.Context, service *service.Service) {
	store := ctx.KVStore(k.storeKey)
	store.Set(service.Hash, k.cdc.MustMarshalBinaryBare(service))
}

// RemoveService removes the entire service metadata struct for given hash.
func (k Keeper) RemoveService(ctx sdk.Context, hash hash.Hash) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(hash)
}

// GetServices retruns all services
func (k Keeper) GetServices(ctx sdk.Context) []*service.Service {
	store := ctx.KVStore(k.storeKey)
	var services []*service.Service
	for it := sdk.KVStorePrefixIterator(store, nil); it.Valid(); it.Next() {
		services = append(services, k.GetService(ctx, hash.Hash(it.Key())))
	}
	return services
}
