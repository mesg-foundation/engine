package xstructhash

import (
	"github.com/cnf/structhash"
)

// Hash takes a data structure and returns a hash string of that data structure
// at the version asked.
//
// This function uses md5 hashing function and default formatter. See also Dump()
// function.
func Hash(c interface{}, version int) string {
	// structhash.Hash always returns nil.
	hash, _ := structhash.Hash(c, version)
	return hash
}
