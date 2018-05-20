package leveldbprotobuf

import (
	"path/filepath"

	"github.com/golang/protobuf/proto"
	"github.com/mesg-foundation/core/config"
	goleveldb "github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type Database struct{}

var (
	driver      *goleveldb.DB
	storagePath string
)

func init() {
	storagePath = filepath.Join(config.ConfigDirectory, "database")
}

func (db *Database) Open() (err error) {
	driver, err = goleveldb.OpenFile(storagePath, nil)
	return
}

func (db *Database) Close() (err error) {
	err = driver.Close()
	return
}

func (db *Database) Insert(collection string, key string, record proto.Message) (err error) {
	bytes, err := proto.Marshal(record)
	if err != nil {
		return
	}
	err = driver.Put(makeKey(collection, key), bytes, nil)
	return
}

func (db *Database) Delete(collection string, key string) (err error) {
	err = driver.Delete(makeKey(collection, key), nil)
	return
}

func (db *Database) Find(collection string, key string, record proto.Message) (err error) {
	bytes, err := driver.Get(makeKey(collection, key), nil)
	if err != nil {
		return
	}
	err = proto.Unmarshal(bytes, record)
	return
}

func (db *Database) All(collection string, new func() proto.Message) (records []proto.Message, err error) {
	iter := driver.NewIterator(util.BytesPrefix(collectionKey(collection)), nil)
	for iter.Next() {
		record := new()
		err = proto.Unmarshal(iter.Value(), record)
		if err != nil {
			return
		}
		records = append(records, record)
	}
	iter.Release()
	err = iter.Error()
	return
}

func (db *Database) Keys(collection string) (keys []string, err error) {
	iter := driver.NewIterator(util.BytesPrefix(collectionKey(collection)), nil)
	for iter.Next() {
		_, key := splitKey(iter.Key())
		keys = append(keys, key)
	}
	iter.Release()
	err = iter.Error()
	return
}
