package xstructhash

import (
	"github.com/cnf/structhash"
)

// Hash takes a data structure and returns a hash of that data structure
// at the version asked.
//
// This function uses sha1 hashing function and default formatter. See also Dump()
// function.
func Hash(c interface{}, version int) []byte {
	return structhash.Sha1(c, version)
}
