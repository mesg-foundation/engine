package serviceapp

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/hash"
	pbtypes "github.com/mesg-foundation/engine/protobuf/types"
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
func (k Keeper) GetService(ctx sdk.Context, hash hash.Hash) *pbtypes.Service {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(hash.String())) {
		return nil
	}

	b := store.Get([]byte(hash.String()))
	var service pbtypes.Service
	k.cdc.MustUnmarshalBinaryBare(b, &service)
	return &service
}

// Sets the entire Whois metadata struct for a name
func (k Keeper) SetService(ctx sdk.Context, service *pbtypes.Service) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(service.Hash), k.cdc.MustMarshalBinaryBare(service))
}

// RemoveService removes the entire service metadata struct for given hash.
func (k Keeper) RemoveService(ctx sdk.Context, hash hash.Hash) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(hash.String()))
}

// GetServices retruns all services
func (k Keeper) GetServices(ctx sdk.Context) []*pbtypes.Service {
	store := ctx.KVStore(k.storeKey)
	var services []*pbtypes.Service
	for it := sdk.KVStorePrefixIterator(store, nil); it.Valid(); it.Next() {
		var service pbtypes.Service
		k.cdc.MustUnmarshalBinaryBare(it.Value(), &service)
		services = append(services, &service)
	}
	return services
}
