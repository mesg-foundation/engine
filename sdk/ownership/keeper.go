package ownershipsdk

import (
	"fmt"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/codec"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/ownership"
)

// Keeper holds the logic to read and write data.
type Keeper struct {
	storeKey *cosmostypes.KVStoreKey
}

// NewKeeper initialize a new keeper.
func NewKeeper(storeKey *cosmostypes.KVStoreKey) *Keeper {
	return &Keeper{
		storeKey: storeKey,
	}
}

// Create creates a new ownership.
func (k *Keeper) Create(req cosmostypes.Request, owner cosmostypes.AccAddress, resourceHash hash.Hash, resource ownership.Ownership_Resource) (*ownership.Ownership, error) {
	store := req.KVStore(k.storeKey)
	hashes := findOwnerships(store, "", resourceHash)
	if len(hashes) > 0 {
		return nil, fmt.Errorf("resource %s:%q already has an owner", resource, resourceHash)
	}
	ownership := &ownership.Ownership{
		Owner:        owner.String(),
		Resource:     resource,
		ResourceHash: resourceHash,
	}
	ownership.Hash = hash.Dump(ownership)

	value, err := codec.MarshalBinaryBare(ownership)
	if err != nil {
		return nil, err
	}
	store.Set(ownership.Hash, value)
	return ownership, nil
}

// Delete deletes a ownership.
func (k *Keeper) Delete(req cosmostypes.Request, owner cosmostypes.AccAddress, resourceHash hash.Hash) error {
	store := req.KVStore(k.storeKey)
	hashes := findOwnerships(store, owner.String(), resourceHash)
	if len(hashes) == 0 {
		return fmt.Errorf("resource %q do not have any ownership", resourceHash)
	}

	for _, hash := range hashes {
		store.Delete(hash)
	}
	return nil
}

// List returns all ownerships.
func (k *Keeper) List(req cosmostypes.Request) ([]*ownership.Ownership, error) {
	var (
		ownerships []*ownership.Ownership
		iter       = req.KVStore(k.storeKey).Iterator(nil, nil)
	)
	defer iter.Close()

	for iter.Valid() {
		var o *ownership.Ownership
		if err := codec.UnmarshalBinaryBare(iter.Value(), &o); err != nil {
			return nil, err
		}
		ownerships = append(ownerships, o)
		iter.Next()
	}
	return ownerships, nil
}

// hasOwner checks if given resource has owner. Returns all ownership hash and true if has one
// nil and false otherwise.
func findOwnerships(store cosmostypes.KVStore, owner string, resourceHash hash.Hash) []hash.Hash {
	var ownerships []hash.Hash
	iter := store.Iterator(nil, nil)

	for iter.Valid() {
		var o *ownership.Ownership
		if err := codec.UnmarshalBinaryBare(iter.Value(), &o); err == nil {
			if (owner == "" || o.Owner == owner) && o.ResourceHash.Equal(resourceHash) {
				ownerships = append(ownerships, o.Hash)
			}
		}
		iter.Next()
	}

	iter.Close()
	return ownerships
}
