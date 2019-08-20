package keeper //TODO: to rename to store

// Store describes the API of keeper of Services.
type Store interface {
	Get(key []byte) ([]byte, error)
	Has(key []byte) (bool, error)
	Delete(key []byte) error
	Put(key []byte, value []byte) error
	NewIterator() Iterator
	Close() error
}

// Iterator
type Iterator interface {
	Next() bool
	Key() []byte
	Value() []byte
	Release()
	Error() error
}
