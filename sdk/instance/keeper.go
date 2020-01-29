package instancesdk

import (
	"fmt"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/codec"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/instance"
	"github.com/mesg-foundation/engine/protobuf/api"
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

// FetchOrCreate creates a new instance if needed.
func (k *Keeper) FetchOrCreate(request cosmostypes.Request, serviceHash hash.Hash, envHash hash.Hash) (*instance.Instance, error) {
	inst := &instance.Instance{
		ServiceHash: serviceHash,
		EnvHash:     envHash,
	}
	inst.Hash = hash.Dump(inst)

	if store := request.KVStore(k.storeKey); !store.Has(inst.Hash) {
		value, err := codec.MarshalBinaryBare(inst)
		if err != nil {
			return nil, err
		}
		store.Set(inst.Hash, value)
	}

	return inst, nil
}

// Get returns the instance that matches given hash.
func (k *Keeper) Get(request cosmostypes.Request, hash hash.Hash) (*instance.Instance, error) {
	var i *instance.Instance
	store := request.KVStore(k.storeKey)
	if !store.Has(hash) {
		return nil, fmt.Errorf("instance %q not found", hash)
	}
	value := store.Get(hash)
	return i, codec.UnmarshalBinaryBare(value, &i)
}

// Exists returns true if a specific set of data exists in the database, false otherwise
func (k *Keeper) Exists(request cosmostypes.Request, hash hash.Hash) (bool, error) {
	return request.KVStore(k.storeKey).Has(hash), nil
}

// List returns all instances.
func (k *Keeper) List(request cosmostypes.Request, f *api.ListInstanceRequest_Filter) ([]*instance.Instance, error) {
	var (
		instances []*instance.Instance
		iter      = request.KVStore(k.storeKey).Iterator(nil, nil)
	)

	// filter results
	for iter.Valid() {
		var i *instance.Instance
		if err := codec.UnmarshalBinaryBare(iter.Value(), &i); err != nil {
			return nil, err
		}
		if f == nil || f.ServiceHash.IsZero() || i.ServiceHash.Equal(f.ServiceHash) {
			instances = append(instances, i)
		}
		iter.Next()
	}
	iter.Close()
	return instances, nil
}
