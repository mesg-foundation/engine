package store

import (
	"github.com/syndtr/goleveldb/leveldb"
)

// LevelDBStore
type LevelDBStore struct {
	db *leveldb.DB
}

func NewLevelDBStore(path string) (*LevelDBStore, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDBStore{
		db: db,
	}, nil
}

func (s *LevelDBStore) NewIterator() Iterator {
	return s.db.NewIterator(nil, nil)
}

func (s *LevelDBStore) Has(key []byte) (bool, error) {
	if _, err := s.db.Get(key, nil); err != nil {
		if err == leveldb.ErrNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *LevelDBStore) Get(key []byte) ([]byte, error) {
	return s.db.Get(key, nil)
}

func (s *LevelDBStore) Delete(key []byte) error {
	return s.db.Delete(key, nil)
}

func (s *LevelDBStore) Put(key []byte, value []byte) error {
	return s.db.Put(key, value, nil)
}

func (s *LevelDBStore) Close() error {
	return s.db.Close()
}
