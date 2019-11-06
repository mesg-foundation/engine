package cosmos

import (
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/hash"
)

// wrapRecoverError recover error from panic and wrap it into a db error.
func wrapRecoverError(err *error) {
	if *err != nil {
		return
	}
	r := recover()
	if r == nil {
		return
	}
	errR, ok := r.(error)
	if !ok {
		errR = fmt.Errorf("%s", r)
	}
	*err = fmt.Errorf("db: %w", errR)
}

// DB is a cosmos store with error and encoding.
type DB struct {
	store  types.KVStore
	cdc *codec.Codec
}

// NewDB returns a new Cosmos store wrapper.
func NewDB(store types.KVStore, cdc *codec.Codec) *DB {
	return &DB{
		store:  store,
		cdc: cdc,
	}
}

// Has returns true if the key is set in the db.
func (db *DB) Has(key hash.Hash) (has bool, err error) {
	defer wrapRecoverError(&err)
	has = db.store.Has(key)
	return
}

// Get retrives the value from db. It returns an error if the db does not contains the key.
func (db *DB) Get(key hash.Hash, ptr interface{}) (err error) {
	defer wrapRecoverError(&err)
	has, err := db.Has(key)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("db: not found")
	}
	return db.cdc.UnmarshalBinaryBare(db.store.Get(key), ptr)
}

// Delete deletes the value for the given key. Delete will not returns error if key doesn't exist.
func (db *DB) Delete(key hash.Hash) (err error) {
	defer wrapRecoverError(&err)
	db.store.Delete(key)
	return
}

// Save sets the value for the given key. It overwrites any previous value.
func (db *DB) Save(key hash.Hash, data interface{}) (err error) {
	defer wrapRecoverError(&err)
	value, err := db.cdc.MarshalBinaryBare(data)
	if err != nil {
		return err
	}
	db.store.Set(key, value)
	return
}

// NewIterator returns a new iterator.
func (db *DB) NewIterator() *DBIterator {
	return &DBIterator{
		iter:  types.KVStorePrefixIterator(db.store, nil),
		valid: false,
		cdc:   db.cdc,
	}
}

// DBIterator is a Cosmos store's iterator implementation of Iterator.
type DBIterator struct {
	iter  types.Iterator
	err   error
	valid bool // HACK for next function. Iterator of cosmos already contains the first element at its creation.
	cdc   *codec.Codec
}

// Next moves the iterator to the next sequential key in the db.
func (i *DBIterator) Next() bool {
	defer i.handleError()
	if i.valid {
		i.iter.Next()
	}
	i.valid = i.iter.Valid()
	return i.valid
}

// Key returns the key of the cursor.
func (i *DBIterator) Key() hash.Hash {
	defer i.handleError()
	return i.iter.Key()
}

// Value returns the data of the cursor.
func (i *DBIterator) Value(ptr interface{}) error {
	defer i.handleError()
	return i.cdc.UnmarshalBinaryBare(i.iter.Value(), ptr)
}

// Release releases the Iterator.
func (i *DBIterator) Release() {
	defer i.handleError()
	i.iter.Close()
}

// Error returns any accumulated error.
func (i *DBIterator) Error() error {
	return i.err
}

// returns any accumulated error.
func (i *DBIterator) handleError() {
	if r := recover(); r != nil {
		errR, ok := r.(error)
		if !ok {
			errR = fmt.Errorf("%s", r)
		}
		i.err = fmt.Errorf("db iterator: %w", errR)
		i.valid = false
	}
}
