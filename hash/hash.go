package hash

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mesg-foundation/engine/hash/serializer"
	"github.com/mr-tron/base58"
)

// size is the default size for the hashing algorithm.
const size = sha256.Size

// sumFunc is the default function to hash.
var sumFunc = sha256.Sum256

var errInvalidLen = errors.New("hash: invalid length")

// A Hash is a type for representing common hash.
type Hash []byte

// Dump takes a structure that implement Serializable and returns its hash.
func Dump(v serializer.Serializable) Hash {
	h := sumFunc([]byte(v.Serialize()))
	return Hash(h[:])
}

// Sum takes a slice of byte and returns its hash.
func Sum(v []byte) Hash {
	h := sumFunc(v)
	return Hash(h[:])
}

// Int returns a new hash from a given integer.
// NOTE: This method is for tests purpose only.
func Int(h int) Hash {
	hash := make(Hash, size)
	binary.PutUvarint(hash, uint64(h))
	return hash
}

// Random returns a new random hash.
func Random() (Hash, error) {
	hash := make(Hash, size)
	n, err := rand.Reader.Read(hash)
	if err != nil {
		return nil, fmt.Errorf("hash generate random error: %s", err)
	}
	if n != size {
		return nil, fmt.Errorf("hash generate random error: invalid hash length")
	}
	return hash, nil
}

// Decode decodes the base58 encoded hash. It returns error
// when a hash isn't base58,encoded or hash length is invalid.
func Decode(h string) (Hash, error) {
	hash, err := base58.Decode(h)
	if err != nil {
		return nil, fmt.Errorf("hash: %s", err)
	}
	if len(hash) != size {
		return nil, errInvalidLen
	}
	return Hash(hash), nil
}

// DecodeFromBytes decodes hash and checks it length.
// It returns empty hash on nil slice of bytes.
func DecodeFromBytes(data []byte) (Hash, error) {
	if len(data) != size {
		return nil, errInvalidLen
	}
	return Hash(data), nil
}

// IsZero reports whethere h represents empty hash.
func (h Hash) IsZero() bool {
	return len(h) == 0
}

// String returns the hash hex representation.
func (h Hash) String() string {
	return base58.Encode(h)
}

// Equal returns a boolean reporting whether h and h1 are the same hashes.
func (h Hash) Equal(h1 Hash) bool {
	return bytes.Equal(h, h1)
}

// Marshal marshals hash into slice of bytes. It's used by protobuf.
func (h Hash) Marshal() ([]byte, error) {
	if len(h) != size {
		return nil, errInvalidLen
	}
	return h, nil
}

// MarshalTo marshals hash into slice of bytes. It's used by protobuf.
func (h Hash) MarshalTo(data []byte) (int, error) {
	if len(h) != size {
		return 0, errInvalidLen
	}
	return copy(data, h), nil
}

// Unmarshal unmarshals slice of bytes into hash. It's used by protobuf.
func (h *Hash) Unmarshal(data []byte) error {
	if len(data) != size {
		return errInvalidLen
	}
	*h = make([]byte, len(data))
	copy(*h, data)
	return nil
}

// Size retruns size of hash. It's used by protobuf.
func (h Hash) Size() int {
	return len(h)
}

// Valid checks if service hash length is valid.
// It treats empty hash as valid one.
func (h Hash) Valid() bool {
	return len(h) == 0 || len(h) == size
}

// MarshalJSON marshals hash into encoded json string.
func (h Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(base58.Encode(h))
}

// UnmarshalJSON unmarshals hex encoded json string into hash.
func (h *Hash) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	if str == "" {
		return nil
	}
	h1, err := base58.Decode(str)
	if err != nil {
		return err
	}
	*h = Hash(h1)
	return nil
}
