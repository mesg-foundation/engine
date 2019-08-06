package serviceapp

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes
// getter/setter methods for the various parts of the service state machine.
type Keeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context
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
func (k Keeper) GetService(ctx sdk.Context, hash string) Service {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(hash)) {
		return Service{}
	}
	b := store.Get([]byte(hash))
	var service Service
	k.cdc.MustUnmarshalBinaryBare(b, &service)
	return service
}

// Sets the entire Whois metadata struct for a name
func (k Keeper) SetService(ctx sdk.Context, service Service) {
	if service.Owner.Empty() {
		return
	}
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(service.Hash), k.cdc.MustMarshalBinaryBare(service))
}

// RemoveService removes the entire service metadata struct for given hash.
func (k Keeper) RemoveService(ctx sdk.Context, hash string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(hash))
}

// GetHashesIterator returns iterator over all hashs in which the keys are the names and the values are the whois
func (k Keeper) GetHashesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}
