package xrand

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"time"
)

// SeedInit initialize rand package with secure random number.
func SeedInit() {
	var b [8]byte
	if _, err := crypto_rand.Read(b[:]); err == nil {
		rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
	} else {
		// fallback to unix nano time
		rand.Seed(time.Now().UTC().UnixNano())
	}
}
