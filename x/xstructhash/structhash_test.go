package xstructhash

import (
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	if got := fmt.Sprintf("%x", Hash(struct{}{}, 1)); got != "bf21a9e8fbc5a3846fb05b4fa0859e0917b2202f" {
		t.Errorf("hash dosen't match - got: %s", got)
	}
}
