package store

import (
	"github.com/syndtr/goleveldb/leveldb"
)

// LevelDBStore is a levelDB implementation of Store.
type LevelDBStore struct {
	db *leveldb.DB
}

// NewLevelDBStore returns a new level db wrapper.
func NewLevelDBStore(path string) (*LevelDBStore, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDBStore{
		db: db,
	}, nil
}

// NewIterator returns a new iterator.
func (s *LevelDBStore) NewIterator() Iterator {
	return s.db.NewIterator(nil, nil)
}

// Has returns true if the key is set in the store.
func (s *LevelDBStore) Has(key []byte) (bool, error) {
	_, err := s.Get(key)
	return err != ErrNotFound, err
}

// Get retrives service from store. It returns ErrNotFound if the store does not contains the key.
func (s *LevelDBStore) Get(key []byte) ([]byte, error) {
	value, err := s.db.Get(key, nil)
	if err == leveldb.ErrNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return value, nil
}

// Delete deletes the value for the given key. Delete will not returns error if key doesn't exist.
func (s *LevelDBStore) Delete(key []byte) error {
	return s.db.Delete(key, nil)
}

// Put sets the value for the given key. It overwrites any previous value.
func (s *LevelDBStore) Put(key []byte, value []byte) error {
	return s.db.Put(key, value, nil)
}

// Close closes the store.
func (s *LevelDBStore) Close() error {
	return s.db.Close()
}
