package database

type Database interface {
	Open() error
	Close() error
	Insert(table string, key string, data interface{}) error
	Find(table string, key string, data interface{}) error
	All(collection string) (data []string, err error)
}
