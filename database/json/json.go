package json

import (
	"path/filepath"

	"github.com/mesg-foundation/core/config"
	"github.com/nanobox-io/golang-scribble"
)

type Database struct{}

var driver *scribble.Driver

var storagePath string

func init() {
	storagePath = filepath.Join(config.ConfigDirectory, "database")
}

func (db *Database) Open() (err error) {
	driver, err = scribble.New(storagePath, nil)
	return
}

func (db *Database) Close() (err error) {
	driver = nil
	return
}

func (db *Database) Insert(collection string, key string, data interface{}) (err error) {
	err = driver.Write(collection, key, data)
	return
}

func (db *Database) Delete(collection string, key string) (err error) {
	err = driver.Delete(collection, key)
	return
}

func (db *Database) Find(collection string, key string, data interface{}) (err error) {
	err = driver.Read(collection, key, &data)
	return
}

func (db *Database) All(collection string) (data []string, err error) {
	data, err = driver.ReadAll(collection)
	return
}
