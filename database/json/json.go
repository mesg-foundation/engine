package json

import (
	"path/filepath"

	"github.com/mesg-foundation/core/config"
	"github.com/nanobox-io/golang-scribble"
)

type Database struct{}

// The
var driver *scribble.Driver

var storagePath string

func init() {
	storagePath = filepath.Join(config.ConfigDirectory, "database")
}

func (db *Database) Open() (err error) {
	driver, err = scribble.New(storagePath, nil)
	return
}

func (db *Database) Insert(table string, key string, data interface{}) (err error) {
	err = driver.Write(table, key, data)
	return
}

func (db *Database) Find(table string, key string, data interface{}) (err error) {
	err = driver.Read(table, key, &data)
	return
}

func (db *Database) Close() (err error) {
	return
}

func (db *Database) All(collection string) (data []string, err error) {
	data, err = driver.ReadAll(collection)
	return
}
