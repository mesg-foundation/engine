package servicesdk

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/hash"
	serv "github.com/mesg-foundation/engine/service"
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

func (k Keeper) Get(ctx sdk.Context, hashString string) serv.Service {
	store := ctx.KVStore(k.storeKey)
	h, err := hash.Decode(hashString)
	if err != nil {
		return serv.Service{}
	}
	if !store.Has(h) {
		return serv.Service{}
	}
	b := store.Get(h)
	var service serv.Service
	k.cdc.MustUnmarshalBinaryBare(b, &service)
	return service
}

// Sets the entire Whois metadata struct for a name
func (k Keeper) Set(ctx sdk.Context, service serv.Service) {
	// if service.Owner.Empty() {
	// 	return
	// }
	store := ctx.KVStore(k.storeKey)
	store.Set(service.Hash, k.cdc.MustMarshalBinaryBare(service))
}

// // RemoveService removes the entire service metadata struct for given hash.
// func (k Keeper) RemoveService(ctx sdk.Context, hash string) {
// 	store := ctx.KVStore(k.storeKey)
// 	store.Delete([]byte(hash))
// }

// // GetHashesIterator returns iterator over all hashs in which the keys are the names and the values are the whois
// func (k Keeper) GetHashesIterator(ctx sdk.Context) sdk.Iterator {
// 	store := ctx.KVStore(k.storeKey)
// 	return sdk.KVStorePrefixIterator(store, nil)
// }
