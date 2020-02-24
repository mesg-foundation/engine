package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/ownership"
	"github.com/mesg-foundation/engine/x/ownership/internal/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the ownership store
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
}

// NewKeeper creates a ownership keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey) Keeper {
	keeper := Keeper{
		storeKey: key,
		cdc:      cdc,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// List returns all ownerships.
func (k *Keeper) List(ctx sdk.Context) ([]*ownership.Ownership, error) {
	var (
		ownerships []*ownership.Ownership
		iter       = ctx.KVStore(k.storeKey).Iterator(nil, nil)
	)
	defer iter.Close()

	for iter.Valid() {
		var o *ownership.Ownership
		if err := k.cdc.UnmarshalBinaryLengthPrefixed(iter.Value(), &o); err != nil {
			return nil, err
		}
		ownerships = append(ownerships, o)
		iter.Next()
	}
	return ownerships, nil
}

// Set creates a new ownership.
func (k Keeper) Set(ctx sdk.Context, owner sdk.AccAddress, resourceHash hash.Hash, resource ownership.Ownership_Resource) (*ownership.Ownership, error) {
	store := ctx.KVStore(k.storeKey)
	hashes := k.findOwnerships(store, "", resourceHash)
	if len(hashes) > 0 {
		return nil, fmt.Errorf("resource %s:%q already has an owner", resource, resourceHash)
	}

	ownership := &ownership.Ownership{
		Owner:        owner.String(),
		Resource:     resource,
		ResourceHash: resourceHash,
	}
	ownership.Hash = hash.Dump(ownership)

	store.Set(ownership.Hash, k.cdc.MustMarshalBinaryLengthPrefixed(ownership))
	return ownership, nil
}

// Delete deletes an ownership.
func (k Keeper) Delete(ctx sdk.Context, owner sdk.AccAddress, resourceHash hash.Hash) error {
	store := ctx.KVStore(k.storeKey)
	hashes := k.findOwnerships(store, owner.String(), resourceHash)
	if len(hashes) == 0 {
		return fmt.Errorf("resource %q do not have any ownership", resourceHash)
	}

	for _, hash := range hashes {
		store.Delete(hash)
	}
	return nil
}

// hasOwner checks if given resource has owner. Returns all ownership hash and true if has one
// nil and false otherwise.
func (k Keeper) findOwnerships(store sdk.KVStore, owner string, resourceHash hash.Hash) []hash.Hash {
	var ownerships []hash.Hash
	iter := store.Iterator(nil, nil)

	for iter.Valid() {
		var o *ownership.Ownership
		if err := k.cdc.UnmarshalBinaryLengthPrefixed(iter.Value(), &o); err == nil {
			if (owner == "" || o.Owner == owner) && o.ResourceHash.Equal(resourceHash) {
				ownerships = append(ownerships, o.Hash)
			}
		}
		iter.Next()
	}

	iter.Close()
	return ownerships
}
