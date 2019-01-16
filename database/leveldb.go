package database

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// txdb is common interface for leveldb database and transaction.
type leveldbTxDB interface {
	Delete(key []byte, wo *opt.WriteOptions) error
	Get(key []byte, ro *opt.ReadOptions) ([]byte, error)
	Has(key []byte, ro *opt.ReadOptions) (bool, error)
	NewIterator(slice *util.Range, ro *opt.ReadOptions) iterator.Iterator
	Put(key, value []byte, wo *opt.WriteOptions) error
	Write(b *leveldb.Batch, wo *opt.WriteOptions) error
}
