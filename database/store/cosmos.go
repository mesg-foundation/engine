package store

import (
	"github.com/cosmos/cosmos-sdk/types"
)

// CosmosStore is a Cosmos KVStore implementation of Store.
type CosmosStore struct {
	store types.KVStore
}

// NewCosmosStore returns a new Cosmos KVStore wrapper.
func NewCosmosStore(store types.KVStore) *CosmosStore {
	return &CosmosStore{
		store: store,
	}
}

// NewIterator returns a new iterator.
func (s *CosmosStore) NewIterator() Iterator {
	return NewCosmosIterator(types.KVStorePrefixIterator(s.store, nil))
}

// Has returns true if the key is set in the store.
func (s *CosmosStore) Has(key []byte) (bool, error) {
	return s.store.Has(key), nil
}

// Get retrives service from store. It returns ErrNotFound if the store does not contains the key.
func (s *CosmosStore) Get(key []byte) ([]byte, error) {
	has, err := s.Has(key)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrNotFound
	}
	return s.store.Get(key), nil
}

// Delete deletes the value for the given key. Delete will not returns error if key doesn't exist.
func (s *CosmosStore) Delete(key []byte) error {
	s.store.Delete(key)
	return nil
}

// Put sets the value for the given key. It overwrites any previous value.
func (s *CosmosStore) Put(key []byte, value []byte) error {
	s.store.Set(key, value)
	return nil
}

// Close closes the store.
func (s *CosmosStore) Close() error {
	return nil
}

// CosmosIterator is a Cosmos KVStore's iterator implementation of Iterator.
type CosmosIterator struct {
	iter types.Iterator
}

// NewCosmosIterator returns a new Cosmos KVStore Iterator wrapper.
func NewCosmosIterator(iter types.Iterator) *CosmosIterator {
	return &CosmosIterator{
		iter: iter,
	}
}

// Next moves the iterator to the next sequential key in the store.
func (i *CosmosIterator) Next() bool {
	if i.iter.Valid() {
		i.iter.Next()
		return true
	}
	return false
}

// Key returns the key of the cursor.
func (i *CosmosIterator) Key() []byte {
	return i.iter.Key()
}

// Value returns the value of the cursor.
func (i *CosmosIterator) Value() []byte {
	return i.iter.Value()
}

// Release releases the Iterator.
func (i *CosmosIterator) Release() {
	i.iter.Close()
}

// Error returns any accumulated error.
func (i *CosmosIterator) Error() error {
	return nil
}
