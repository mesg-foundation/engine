package store

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

// LevelDBStore
type LevelDBStore struct {
	db *leveldb.DB
}

func NewLevelDBStore(db *leveldb.DB) *LevelDBStore {
	return &LevelDBStore{
		db: db,
	}
}

func (s *LevelDBStore) NewIterator() Iterator {
	return NewLevelDBIterator(s.db.NewIterator(nil, nil))
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

// LevelDBIterator
type LevelDBIterator struct {
	iter iterator.Iterator
}

func NewLevelDBIterator(iter iterator.Iterator) *LevelDBIterator {
	return &LevelDBIterator{
		iter: iter,
	}
}

func (i *LevelDBIterator) Next() bool {
	return i.iter.Next()
}

func (i *LevelDBIterator) Key() []byte {
	return i.iter.Key()
}

func (i *LevelDBIterator) Value() []byte {
	return i.iter.Value()
}

func (i *LevelDBIterator) Release() {
	i.iter.Release()
}

func (i *LevelDBIterator) Error() error {
	return i.iter.Error()
}
