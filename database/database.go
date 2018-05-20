package database

type Record interface {
	Key() (key string)
	Encode() (bytes []byte, err error)
	Decode(bytes []byte) (err error)
}

type Database interface {
	Open() (err error)
	Close() (err error)
	Insert(collection string, key string, record Record) (err error)
	Delete(collection string, key string) (err error)
	Find(collection string, key string, record Record) (err error)
	All(collection string, new func() Record) (records []Record, err error)
}

var Db Database

func init() {
	// Db = &json.Database{}
}
