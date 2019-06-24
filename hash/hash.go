package hash

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash"

	"github.com/cnf/structhash"
	"github.com/mr-tron/base58"
)

// DefaultHash is a default hashing algorithm - "sha256".
var DefaultHash = sha256.New

// size is a default size for hashing algorithm.
var size = DefaultHash().Size()

// Digest represents the partial evaluation of a checksum.
type Digest struct {
	hash.Hash
}

// Sum appends the current checksum to b and returns the Hash.
func (d *Digest) Sum(b []byte) Hash {
	return Hash(d.Hash.Sum(b))
}

// A Hash is a type for representing hash with base58 encode and decode functions.
type Hash []byte

// New returns new hash from a given integer.
func New() *Digest {
	return &Digest{
		Hash: DefaultHash(),
	}
}

// Dump takes an interface and returns its hash representation.
func Dump(v interface{}) Hash {
	d := New()
	d.Write(structhash.Dump(v, 0))
	return d.Sum(nil)
}

// Int returns a new hash from a given integer.
// NOTE: This method is for tests purpose only.
func Int(h int) Hash {
	hash := make(Hash, size)
	binary.PutUvarint(hash, uint64(h))
	return hash
}

// Decode decodes the base58 encoded hash. It returns error
// when a hash isn't base58,encoded or hash length is invalid.
func Decode(h string) (Hash, error) {
	hash, err := base58.Decode(h)
	if err != nil {
		return nil, fmt.Errorf("hash: %s", err)
	}
	if len(hash) != size {
		return nil, fmt.Errorf("hash: invalid length string")
	}
	return Hash(hash), nil
}

// IsZero reports whethere h represents empty hash.
func (h Hash) IsZero() bool {
	return len(h) == 0
}

// String returns the base58 hash representation.
func (h Hash) String() string {
	return base58.Encode(h)
}

// Equal returns a boolean reporting whether h and h1 are the same hashes.
func (h Hash) Equal(h1 Hash) bool {
	return bytes.Equal(h, h1)
}
