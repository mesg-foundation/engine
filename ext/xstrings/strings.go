package xstrings

import (
	"math/rand"

	"github.com/mesg-foundation/engine/ext/xrand"
)

func init() {
	xrand.SeedInit()
}

// SliceContains returns true if slice a contains e element, false otherwise.
func SliceContains(a []string, e string) bool {
	for _, s := range a {
		if s == e {
			return true
		}
	}
	return false
}

// FindLongest finds the length of longest string in slice.
func FindLongest(ss []string) int {
	l := 0
	for _, s := range ss {
		if i := len(s); i > l {
			l = i
		}
	}
	return l
}

// SliceIndex returns the index e in a, return -1 if not found.
func SliceIndex(a []string, e string) int {
	for i, s := range a {
		if s == e {
			return i
		}
	}
	return -1
}

var asciiletters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandASCIILetters generates random string from ascii letters.
func RandASCIILetters(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = asciiletters[rand.Intn(len(asciiletters))]
	}
	return string(b)
}
