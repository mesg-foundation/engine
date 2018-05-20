package leveldb

import (
	"encoding/json"
	"path/filepath"

	"github.com/mesg-foundation/core/config"
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

func (db *Database) Insert(collection string, key string, data interface{}) (err error) {
	bin, err := json.Marshal(data)
	err = driver.Put(makeKey(collection, key), bin, nil)
	return
}

func (db *Database) Delete(collection string, key string) (err error) {
	err = driver.Delete(makeKey(collection, key), nil)
	return
}

func (db *Database) Find(collection string, key string, data interface{}) (err error) {
	bin, err := driver.Get(makeKey(collection, key), nil)
	if err != nil {
		return
	}
	err = json.Unmarshal(bin, &data)
	return
}

func (db *Database) All(collection string) (data [][]byte, err error) {
	iter := driver.NewIterator(util.BytesPrefix([]byte(collection)), nil)
	for iter.Next() {
		value := iter.Value()
		dest := make([]byte, len(value))
		copy(dest, value)
		data = append(data, dest)
	}
	iter.Release()
	err = iter.Error()
	return
}
