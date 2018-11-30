package xstructhash

import (
	"crypto/sha1"
	"fmt"
	"testing"

	"github.com/cnf/structhash"
)

func TestHash(t *testing.T) {
	s := struct{}{}
	v := 1
	got := Hash(s, v)
	want := fmt.Sprintf("%x", sha1.Sum(structhash.Dump(s, v)))
	if got != want {
		t.Errorf("invalid hash")
	}
}
