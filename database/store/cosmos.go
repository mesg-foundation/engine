package store

import (
	"errors"
	"fmt"

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
func (s *CosmosStore) Has(key []byte) (has bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			if err, ok = r.(error); !ok {
				err = fmt.Errorf("store: %s", r)
			}
		}
	}()
	has = s.store.Has(key)
	return
}

// Get retrives service from store. It returns an error if the store does not contains the key.
func (s *CosmosStore) Get(key []byte) (out []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			if err, ok = r.(error); !ok {
				err = fmt.Errorf("store: %s", r)
			}
		}
	}()

	has, err := s.Has(key)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("not found")
	}
	out = s.store.Get(key)
	return
}

// Delete deletes the value for the given key. Delete will not returns error if key doesn't exist.
func (s *CosmosStore) Delete(key []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			if err, ok = r.(error); !ok {
				err = fmt.Errorf("store: %s", r)
			}
		}
	}()
	s.store.Delete(key)
	return
}

// Put sets the value for the given key. It overwrites any previous value.
func (s *CosmosStore) Put(key []byte, value []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			if err, ok = r.(error); !ok {
				err = fmt.Errorf("store: %s", r)
			}
		}
	}()
	s.store.Set(key, value)
	return
}

// Close closes the store.
func (s *CosmosStore) Close() error {
	return nil
}

// CosmosIterator is a Cosmos KVStore's iterator implementation of Iterator.
type CosmosIterator struct {
	iter  types.Iterator
	err   error
	valid bool // HACK for next function. Iterator of cosmos already contains the first element at its creation.
}

// NewCosmosIterator returns a new Cosmos KVStore Iterator wrapper.
func NewCosmosIterator(iter types.Iterator) *CosmosIterator {
	return &CosmosIterator{
		iter:  iter,
		valid: false,
	}
}

// Next moves the iterator to the next sequential key in the store.
func (i *CosmosIterator) Next() bool {
	defer i.handleError()
	if i.valid {
		i.iter.Next()
	}
	i.valid = i.iter.Valid()
	return i.valid
}

// Key returns the key of the cursor.
func (i *CosmosIterator) Key() []byte {
	defer i.handleError()
	return i.iter.Key()
}

// Value returns the value of the cursor.
func (i *CosmosIterator) Value() []byte {
	defer i.handleError()
	return i.iter.Value()
}

// Release releases the Iterator.
func (i *CosmosIterator) Release() {
	defer i.handleError()
	i.iter.Close()
}

// Error returns any accumulated error.
func (i *CosmosIterator) Error() error {
	return i.err
}

// returns any accumulated error.
func (i *CosmosIterator) handleError() {
	if r := recover(); r != nil {
		var ok bool
		if i.err, ok = r.(error); !ok {
			i.err = fmt.Errorf("store iterator: %s", r)
		}
	}
}
