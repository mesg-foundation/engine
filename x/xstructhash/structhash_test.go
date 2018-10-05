package xstructhash

import (
	"testing"

	"github.com/cnf/structhash"
)

func TestHash(t *testing.T) {
	hash := Hash(struct{}{}, 1)
	if version := structhash.Version(hash); version != 1 {
		t.Errorf("invalid version - got %d, want %d", version, 1)
	}
}
