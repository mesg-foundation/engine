package xstructhash

import (
	"fmt"

	"github.com/cnf/structhash"
)

// Hash takes a data structure and returns a hash string of that data structure
// at the version asked.
//
// This function uses md5 hashing function and default formatter. See also Dump()
// function.
func Hash(c interface{}, version int) string {
	return fmt.Sprintf("%x", structhash.Sha1(c, version))
}
