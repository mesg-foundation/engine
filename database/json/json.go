package json

import (
	"path/filepath"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/database"
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

func (db *Database) Insert(collection string, key string, record database.Record) (err error) {
	err = driver.Write(collection, key, record)
	return
}

func (db *Database) Delete(collection string, key string) (err error) {
	err = driver.Delete(collection, key)
	return
}

func (db *Database) Find(collection string, key string, record database.Record) (err error) {
	err = driver.Read(collection, key, &record)
	return
}

func (db *Database) All(collection string, new func() database.Record) (records []database.Record, err error) {
	strings, err := driver.ReadAll(collection)
	records = make([]database.Record, len(strings))
	for i, element := range strings {
		records[i] = new()
		err = records[i].Decode([]byte(element))
		if err != nil {
			return
		}
	}
	return
}
