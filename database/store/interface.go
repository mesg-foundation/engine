package store

// Store describes the public API of a store.
type Store interface {
	// Get gets the value for the given key. It returns nil if the store does not contains the key.
	Get(key []byte) ([]byte, error)

	// Has returns true if the key is set in the store.
	Has(key []byte) (bool, error)

	// Delete deletes the value for the given key. Delete will not returns error if key doesn't exist.
	Delete(key []byte) error

	// Put sets the value for the given key. It overwrites any previous value.
	Put(key []byte, value []byte) error

	// NewIterator returns a new iterator.
	NewIterator() Iterator

	// Close closes the store.
	Close() error
}

// Iterator describes the public API of an iterator.
type Iterator interface {
	// Next moves the iterator to the next sequential key in the database.
	Next() bool

	// Key returns the key of the cursor.
	Key() []byte

	// Value returns the value of the cursor.
	Value() []byte

	// Release releases the Iterator.
	Release()

	// Error returns any accumulated error.
	Error() error
}
