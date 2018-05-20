package leveldb

import (
	"path/filepath"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/database"
	goleveldb "github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type Database struct{}

var driver *goleveldb.DB

var storagePath string

func init() {
	storagePath = filepath.Join(config.ConfigDirectory, "database")
}

func makeKey(collection string, key string) []byte {
	return []byte(collection + "_" + key)
}

func (db *Database) Open() (err error) {
	driver, err = goleveldb.OpenFile(storagePath, nil)
	return
}

func (db *Database) Close() (err error) {
	driver.Close()
	return
}

func (db *Database) Insert(collection string, key string, record database.Record) (err error) {
	bytes, err := record.Encode()
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

func (db *Database) Find(collection string, key string, record database.Record) (err error) {
	bytes, err := driver.Get(makeKey(collection, key), nil)
	if err != nil {
		return
	}
	err = record.Decode(bytes)
	return
}

func (db *Database) All(collection string, new func() database.Record) (records []database.Record, err error) {
	iter := driver.NewIterator(util.BytesPrefix([]byte(collection)), nil)
	for iter.Next() {
		record := new()
		err = record.Decode(iter.Value())
		if err != nil {
			return
		}
		records = append(records, record)
	}
	iter.Release()
	err = iter.Error()
	return
}
