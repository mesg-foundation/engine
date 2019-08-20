package keeper

import (
	"github.com/cosmos/cosmos-sdk/types"
)

// CosmosStore
type CosmosStore struct {
	store types.KVStore
}

func NewCosmosStore(store types.KVStore) *CosmosStore {
	return &CosmosStore{
		store: store,
	}
}

func (s *CosmosStore) NewIterator() Iterator {
	return NewCosmosIterator(types.KVStorePrefixIterator(s.store, nil))
}

func (s *CosmosStore) Has(key []byte) (bool, error) {
	return s.store.Has(key), nil
}

func (s *CosmosStore) Get(key []byte) ([]byte, error) {
	return s.store.Get(key), nil
}

func (s *CosmosStore) Delete(key []byte) error {
	s.store.Delete(key)
	return nil
}

func (s *CosmosStore) Put(key []byte, value []byte) error {
	s.store.Set(key, value)
	return nil
}

func (s *CosmosStore) Close() error {
	return nil
}

// CosmosIterator
type CosmosIterator struct {
	iter types.Iterator
}

func NewCosmosIterator(iter types.Iterator) *CosmosIterator {
	return &CosmosIterator{
		iter: iter,
	}
}

func (i *CosmosIterator) Next() bool {
	i.iter.Next()
	return i.iter.Valid()
}

func (i *CosmosIterator) Key() []byte {
	return i.iter.Key()
}

func (i *CosmosIterator) Value() []byte {
	return i.iter.Value()
}

func (i *CosmosIterator) Release() {
	i.iter.Close()
}

func (i *CosmosIterator) Error() error {
	return nil
}
