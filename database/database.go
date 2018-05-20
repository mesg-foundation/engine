package database

import "github.com/mesg-foundation/core/database/json"

type Database interface {
	Open() (err error)
	Close() (err error)
	Insert(collection string, key string, data interface{}) (err error)
	Delete(collection string, key string) (err error)
	Find(collection string, key string, data interface{}) (err error)
	All(collection string) (data [][]byte, err error)
}

var Db Database

func init() {
	Db = &json.Database{}
}
